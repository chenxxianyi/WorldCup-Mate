// @vitest-environment jsdom
import { describe, expect, it, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import TeamFlag from '@/components/common/TeamFlag.vue'

describe('TeamFlag', () => {
  it('renders an image when value is an http URL', () => {
    const wrapper = mount(TeamFlag, { props: { value: 'https://cdn.example/flag.png', alt: 'ARG' } })
    const img = wrapper.find('img')
    expect(img.exists()).toBe(true)
    expect(img.attributes('src')).toBe('https://cdn.example/flag.png')
  })

  it('falls back to the flag-text span for non-image values', () => {
    // "XYZ" has no local team icon and is not an image URL.
    const wrapper = mount(TeamFlag, { props: { value: '测试队', fallback: 'XYZ' } })
    expect(wrapper.find('img').exists()).toBe(false)
    expect(wrapper.find('.flag-text').text()).toContain('测试队')
  })

  it('shows the fallback chip when the image fails to load', async () => {
    const wrapper = mount(TeamFlag, { props: { value: 'https://cdn.example/broken.png', alt: 'BRA' } })
    await wrapper.find('img').trigger('error')
    expect(wrapper.find('img').exists()).toBe(false)
    // ARG-style art code is not in the art set for this value, so the
    // generic fallback chip shows.
    expect(wrapper.find('.flag-fallback').text()).toBe('BRA')
  })

  it('re-tries the image when the value changes after a failure', async () => {
    const wrapper = mount(TeamFlag, { props: { value: 'https://cdn.example/broken.png', alt: 'BRA' } })
    await wrapper.find('img').trigger('error')
    await wrapper.setProps({ value: 'https://cdn.example/ok.png' })
    expect(wrapper.find('img').exists()).toBe(true)
    expect(wrapper.find('img').attributes('src')).toBe('https://cdn.example/ok.png')
  })
})
