// @vitest-environment jsdom
import { describe, expect, it } from 'vitest'
import { mount } from '@vue/test-utils'
import Countdown from '@/components/common/Countdown.vue'

describe('Countdown', () => {
  it('renders hh:mm:ss from explicit parts', async () => {
    const wrapper = mount(Countdown, { props: { hours: 1, minutes: 5, seconds: 9 } })
    // onMounted recomputes the display; flush the reactive update.
    await wrapper.vm.$nextTick()
    const boxes = wrapper.findAll('.time-box strong')
    expect(boxes.map((b) => b.text())).toEqual(['01', '05', '09'])
  })

  it('renders zeros for a zero countdown', async () => {
    const wrapper = mount(Countdown, { props: { seconds: 0 } })
    await wrapper.vm.$nextTick()
    expect(wrapper.findAll('.time-box strong').map((b) => b.text())).toEqual(['00', '00', '00'])
  })
})
