<script setup lang="ts">
import { Result } from '../../bindings/github.com/maicek/laluer/core/handler';

const { data } = defineProps<{
  active: boolean;
  data: Result;
}>();
</script>

<template>
  <div class="Item" :class="{ active: active }">
    <img
      class="Item__icon"
      :src="`data:image/svg+xml;base64,${data.iconBase64}`"
      :alt="data.label"
      @error="
        (item) => (item.target as HTMLImageElement).classList.add('error')
      "
    />

    <div class="Item__content">
      <span class="Item__content-label">{{ data.label }}</span>
      <span class="Item__content-description">{{ data.subtitle }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.Item {
  border-radius: 10px;
  padding: 10px;

  border: 2px solid #cccccc40;
  background-color: rgba(15, 26, 46, 0.2);

  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;

  transition: all 0.2s ease;

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
    width: 32px;
    height: 32px;
    object-fit: contain;

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
