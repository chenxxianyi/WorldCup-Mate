<script setup lang="ts">
defineProps<{
  modelValue: string
  placeholder?: string
}>()

defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()
</script>

<template>
  <div class="search-wrap">
    <input
      class="search"
      :value="modelValue"
      :placeholder="placeholder || '搜索'"
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <button
      v-if="modelValue"
      class="clear-btn"
      type="button"
      aria-label="清空搜索"
      title="清空搜索"
      @click="$emit('update:modelValue', '')"
    >
      <span class="material-symbols-outlined">close</span>
    </button>
  </div>
</template>

<style scoped>
.search-wrap {
  position: relative;
  width: 100%;
}

.search {
  width: 100%;
  height: 46px;
  border: 1px solid var(--line);
  border-radius: 16px;
  outline: none;
  padding: 0 48px 0 15px;
  color: var(--text);
  background: var(--card);
  box-shadow: var(--shadow);
  font-size: 14px;
}

.search::placeholder {
  color: var(--weak);
}

.clear-btn {
  position: absolute;
  top: 50%;
  right: 10px;
  width: 30px;
  height: 30px;
  display: grid;
  place-items: center;
  border: 0;
  border-radius: 50%;
  color: var(--muted);
  background: transparent;
  transform: translateY(-50%);
}

.clear-btn:active {
  color: var(--primary);
  background: color-mix(in srgb, var(--primary) 8%, transparent);
}

.clear-btn .material-symbols-outlined {
  font-size: 20px;
  line-height: 1;
}
</style>
