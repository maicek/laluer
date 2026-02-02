<script setup lang="ts">
import { reactive, ref, watch } from 'vue';
import {
  Handle,
  Call,
} from '../bindings/github.com/maicek/laluer/core/handler/handlerservice';
import {
  HandlerResult,
  SearchParams,
} from '../bindings/github.com/maicek/laluer/core/handler';
import {} from '../bindings/github.com/maicek/laluer/core/handler';
import Item from './components/Item.vue';
import { useEventListener } from '@vueuse/core';

const searchQuery = ref('');
const results = ref<HandlerResult>({ items: [] });
const activeIndex = ref(0);

watch(searchQuery, async () => {
  results.value = await Handle({
    query: searchQuery.value,
  });
  activeIndex.value = 0;
});

useEventListener('keydown', (e) => {
  switch (e.key) {
    case 'ArrowUp':
      activeIndex.value = Math.max(0, activeIndex.value - 1);
      e.preventDefault();
      break;
    case 'ArrowDown':
      activeIndex.value = Math.min(
        results.value.items.length - 1,
        activeIndex.value + 1
      );
      e.preventDefault();
      break;
    case 'Enter':
      if (results.value.items.length > 0) {
        Call(results.value.items[activeIndex.value].action);
      }
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
    <input v-model="searchQuery" keydown.arrow.up.prevent="" autofocus="true" />

    <div class="Results">
      <template v-for="(item, index) in results.items">
        <Item :data="item" :active="index === activeIndex"></Item>
      </template>
    </div>
  </div>
</template>

<style scoped>
.App {
  width: 100vw;
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 10px 10px;
  gap: 10px;
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);

  input {
    width: 100%;
    height: 40px;
    font-size: 20px;
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 4px;
    text-align: center;
  }
}

.Results {
  display: flex;
  flex-direction: column;
  gap: 5px;
  overflow: auto;
  flex: 1 1 auto;
}
</style>
