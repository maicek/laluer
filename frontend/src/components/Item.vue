<script setup lang="ts">
import { computed } from 'vue';
import { Result } from '../../bindings/github.com/maicek/laluer/core/handler';

const props = defineProps<{
  active: boolean;
  data: Result;
}>();

const iconSrc = computed(() => {
  if (!props.data.iconBase64) return '';
  const mime = props.data.iconMime || 'image/svg+xml';
  return `data:${mime};base64,${props.data.iconBase64}`;
});

const hasIcon = computed(() => !!props.data.iconBase64);

const subtitle = computed(() => {
  return props.data.subtitle || '';
});
</script>

<template>
  <div class="Item" :class="{ active: active }">
    <div class="Item__icon-wrap">
      <img
        v-if="hasIcon"
        class="Item__icon"
        :src="iconSrc"
        :alt="data.label"
        loading="lazy"
      />
      <div v-else class="Item__icon-fallback">
        {{ data.label?.charAt(0)?.toUpperCase() }}
      </div>
    </div>

    <div class="Item__text">
      <span class="Item__label">{{ data.label }}</span>
      <span v-if="subtitle" class="Item__subtitle">{{ subtitle }}</span>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.Item {
  background: rgba(255, 255, 255, 0.04);
  border-radius: 8px;
  padding: 8px 12px;
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  transition: background 0.1s ease;
  min-height: 44px;

  &:hover {
    background: rgba(255, 255, 255, 0.08);
  }
}

.Item.active {
  background: rgba(100, 120, 255, 0.25);
  box-shadow: inset 0 0 0 1px rgba(100, 120, 255, 0.3);
}

.Item__icon-wrap {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.Item__icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
  border-radius: 4px;
}

.Item__icon-fallback {
  width: 32px;
  height: 32px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.5);
}

.Item__text {
  display: flex;
  flex-direction: column;
  min-width: 0;
  gap: 1px;
}

.Item__label {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.9);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.Item__subtitle {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
