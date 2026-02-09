<script setup lang="ts">
import { ref, useTemplateRef, watch } from 'vue';
import { Result } from '../../bindings/github.com/maicek/laluer/core/handler';

const props = defineProps<{
  active: boolean;
  data: Result;
}>();
</script>

<template>
  <div class="Item" :class="{ active: active }">
    <img
      class="Item__icon"
      :src="props.data.iconBase64"
      :alt="props.data.label"
      :key="props.data.iconBase64"
      @error="
        (item) => (item.target as HTMLImageElement).classList.add('error')
      "
    />

    <div class="Item__content">
      <span class="Item__content-label">{{ props.data.label }}</span>
      <span class="Item__content-description">{{ props.data.subtitle }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.Item {
  flex-shrink: 0;
  border-radius: 10px;
  padding: 10px;

  border: 2px solid #cccccc40;
  background-color: rgba(15, 26, 46, 0.2);

  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
  height: 60px !important;
  // max-height: 60px;
  overflow: hidden;

  // transition: all 0.03s ease;

  .Item__content {
    display: flex;
    flex-direction: column;
    gap: 5px;
    line-height: 100%;

    .Item__content-label {
      font-weight: bold;
    }

    .Item__content-description {
      font-size: 12px;
      color: #cccccc80;
      font-weight: 300;
    }
  }

  img {
    width: 32px !important;
    height: 32px !important;
    object-fit: contain;
    aspect-ratio: 1/1 !important;

    &.error {
      display: none;
    }
  }
}

.Item.active {
  background-color: rgba(45, 45, 180, 0.5);
  border-color: #cccccc60;
}
</style>
