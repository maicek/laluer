<script setup lang="ts">
import { ref, watch } from 'vue';
import { Result } from '../../bindings/github.com/maicek/laluer/core/handler';

const props = defineProps<{
  active: boolean;
  data: Result;
}>();

const itemRef = ref<HTMLElement | null>(null);

watch(() => props.active, (isActive) => {
  if (isActive && itemRef.value) {
    itemRef.value.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
  }
});
</script>

<template>
  <div
    ref="itemRef"
    class="Item"
    v-if="props.data.iconBase64 !== '' && props.data.iconBase64 !== null"
    :style="{ display: props.data.iconBase64 !== '' && props.data.iconBase64 !== null ? 'flex' : 'none' }"
    :class="{ active: active }"
  >
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
