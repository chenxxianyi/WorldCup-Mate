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
  <input
    class="search"
    :value="modelValue"
    :placeholder="placeholder || '搜索'"
    @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
  />
</template>

<style scoped>
.search {
  width: 100%;
  height: 46px;
  border: 1px solid var(--line);
  border-radius: 16px;
  outline: none;
  padding: 0 15px;
  color: var(--text);
  background: var(--card);
  box-shadow: var(--shadow);
  font-size: 14px;
}

.search::placeholder {
  color: var(--weak);
}
</style>
