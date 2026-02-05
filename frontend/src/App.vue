<script setup lang="ts">
import { ref, watch, nextTick } from 'vue';
import {
  Handle,
  Call,
} from '../bindings/github.com/maicek/laluer/core/handler/handlerservice';
import {
  HandlerResult,
} from '../bindings/github.com/maicek/laluer/core/handler';
import { Application } from '@wailsio/runtime';
import Item from './components/Item.vue';
import { useEventListener } from '@vueuse/core';

const MAX_VISIBLE = 7;

const searchQuery = ref('');
const results = ref<HandlerResult>({ items: [] });
const activeIndex = ref(0);
const resultsContainer = ref<HTMLElement | null>(null);

const visibleItems = () => {
  return results.value.items.slice(0, MAX_VISIBLE);
};

watch(searchQuery, async () => {
  results.value = await Handle({
    query: searchQuery.value,
  });
  activeIndex.value = 0;
});

function scrollActiveIntoView() {
  nextTick(() => {
    const container = resultsContainer.value;
    if (!container) return;
    const activeEl = container.querySelector('.Item.active') as HTMLElement;
    if (activeEl) {
      activeEl.scrollIntoView({ block: 'nearest', behavior: 'smooth' });
    }
  });
}

useEventListener('keydown', (e) => {
  const items = visibleItems();
  switch (e.key) {
    case 'ArrowUp':
      activeIndex.value = Math.max(0, activeIndex.value - 1);
      scrollActiveIntoView();
      e.preventDefault();
      break;
    case 'ArrowDown':
      activeIndex.value = Math.min(
        items.length - 1,
        activeIndex.value + 1
      );
      scrollActiveIntoView();
      e.preventDefault();
      break;
    case 'Enter':
      if (items.length > 0) {
        Call(items[activeIndex.value].action);
      }
      break;
    case 'Escape':
      Application.Quit();
      break;
    case 'ArrowLeft':
    case 'ArrowRight':
      e.preventDefault();
      break;
  }
});
</script>

<template>
  <div class="App">
    <div class="SearchBar">
      <svg class="SearchIcon" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
        <circle cx="11" cy="11" r="8"/>
        <path d="m21 21-4.3-4.3"/>
      </svg>
      <input
        v-model="searchQuery"
        placeholder="Search applications..."
        autofocus
      />
    </div>

    <div class="Results" ref="resultsContainer">
      <TransitionGroup name="list">
        <Item
          v-for="(item, index) in visibleItems()"
          :key="item.label + item.icon"
          :data="item"
          :active="index === activeIndex"
        />
      </TransitionGroup>
    </div>

    <div class="Hint" v-if="visibleItems().length > 0">
      <span class="Key">&#8593;&#8595;</span> navigate
      <span class="Key">&#9166;</span> launch
      <span class="Key">esc</span> close
    </div>
  </div>
</template>

<style scoped>
.App {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 8px;
  gap: 6px;
  background: rgba(20, 20, 30, 0.85);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  overflow: hidden;
}

.SearchBar {
  display: flex;
  align-items: center;
  gap: 10px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  padding: 0 14px;
  flex-shrink: 0;
}

.SearchIcon {
  width: 18px;
  height: 18px;
  color: rgba(255, 255, 255, 0.4);
  flex-shrink: 0;
}

.SearchBar input {
  width: 100%;
  height: 42px;
  font-size: 16px;
  background: transparent;
  border: none;
  outline: none;
  color: rgba(255, 255, 255, 0.95);
  font-family: inherit;
}

.SearchBar input::placeholder {
  color: rgba(255, 255, 255, 0.3);
}

.Results {
  display: flex;
  flex-direction: column;
  gap: 3px;
  overflow-y: auto;
  overflow-x: hidden;
  flex: 1 1 auto;
  scrollbar-width: none;
}

.Results::-webkit-scrollbar {
  display: none;
}

.Hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 4px 0;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.25);
  flex-shrink: 0;
}

.Hint .Key {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
  padding: 1px 5px;
  font-size: 10px;
  margin-right: 3px;
}

/* TransitionGroup animations */
.list-enter-active {
  transition: all 0.15s ease-out;
}
.list-leave-active {
  transition: all 0.1s ease-in;
}
.list-enter-from {
  opacity: 0;
  transform: translateY(-8px);
}
.list-leave-to {
  opacity: 0;
}
.list-move {
  transition: transform 0.15s ease;
}
</style>
