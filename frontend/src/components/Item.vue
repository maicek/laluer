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
    />

    {{ data.label }}
  </div>
</template>

<style lang="scss" scoped>
.Item {
  background-color: #ffffff30;
  border-radius: 10px;
  padding: 10px;

  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;

  img {
    width: 32px;
    height: 32px;
    object-fit: contain;
  }
}

.Item.active {
  background-color: rgba(45, 45, 180, 0.5);
}
</style>
