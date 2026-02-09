<script setup lang="ts">
import { reactive, ref, useTemplateRef, watch, watchEffect } from 'vue';
import {
  Handle,
  Call,
} from '../bindings/github.com/maicek/laluer/core/handler/handlerservice';
import {
  HandlerResult,
  SearchParams,
} from '../bindings/github.com/maicek/laluer/core/handler';
import { Application } from '@wailsio/runtime';
import Item from './components/Item.vue';
import { useEventListener } from '@vueuse/core';

const searchQuery = ref('');
const parent = useTemplateRef('parent');
const results = ref<HandlerResult>({ items: [] });
const activeIndex = ref(0);

watch(
  searchQuery,
  async () => {
    results.value = await Handle({
      query: searchQuery.value,
    });
    activeIndex.value = 0;
  },
  {
    immediate: true,
  },
);

watchEffect(() => {
  if (activeIndex.value < 0) {
    activeIndex.value = 0;
  } else if (activeIndex.value > results.value.items.length - 1) {
    activeIndex.value = results.value.items.length - 1;
  }
});

useEventListener('keydown', (e) => {
  switch (e.key) {
    case 'Escape':
      Application.Quit();
      break;
    case 'ArrowUp':
      activeIndex.value = Math.max(0, activeIndex.value - 1);
      e.preventDefault();
      break;
    case 'ArrowDown':
      activeIndex.value = Math.min(
        results.value.items.length - 1,
        activeIndex.value + 1,
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

const handleWheelInput = (e: WheelEvent) => {
  e.preventDefault();

  e.deltaY > 0 ? activeIndex.value++ : activeIndex.value--;
  activeIndex.value = Math.max(
    0,
    Math.min(activeIndex.value, results.value.items.length - 1),
  );
};

const items = useTemplateRef('items');

watchEffect(() => {
  const rect = items.value?.at(activeIndex.value)?.$el;

  parent.value?.scrollTo({
    top: rect?.offsetTop - parent.value?.offsetTop - 130,
  });
});

const nekoImageUrl = ref('')

const fetchNekoImage = async () => {
    fetch('https://api.nekosia.cat/api/v1/images/catgirl')
        .then(res => res.json())
        .then(json => {
            nekoImageUrl.value = json.image.original.url
        })
}

fetchNekoImage()

</script>

<template>
  <div class="App">
    <input v-model="searchQuery" keydown.arrow.up.prevent="" autofocus="true" />

    <div class="Results" @wheel="handleWheelInput" ref="parent">
      <template v-for="(item, index) in results.items">
        <Item :data="item" :active="index === activeIndex" ref="items"> </Item>
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

  overflow: hidden !important;

  ::before {
    content: "";
    position: absolute;
    inset: 0;

    background-image: v-bind('`url("${nekoImageUrl}")`');
    background-position: center top;
    background-size: cover;
    background-repeat: no-repeat;

    filter: blur(2px);
    transform: scale(1.1);

    z-index: -2;


  }

  &::after {
    content: "";
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: -1;

  }

  input {
    box-sizing: border-box;
    outline: none;
    width: 100%;
    height: 50px;
    font-size: 20px;
    padding: 10px;
    border: 2px solid rgba(32, 64, 122, 0.6);
    background-color: rgba(33, 33, 33, 0.44);
    border-radius: 8px;
    text-align: center;
    flex: 0 0 50px;

    &:focus {
      outline: none;
    }
  }
}



.Results {
  display: flex;
  flex-direction: column;
  gap: 5px;
  overflow: auto;
  flex: 1 1 auto;
  overflow: hidden;

  & > * {
    box-sizing: border-box;
  }
}
</style>
