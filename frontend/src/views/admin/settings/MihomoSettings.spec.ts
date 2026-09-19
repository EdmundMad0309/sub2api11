import { mount, flushPromises } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import MihomoSettings from './MihomoSettings.vue'
const { get, post } = vi.hoisted(() => ({ get: vi.fn(), post: vi.fn() }))
vi.mock('@/api/client', () => ({ apiClient: { get, post } }))
vi.mock('vue-i18n', () => ({ useI18n: () => ({ locale: { value: 'zh' } }) }))
const base = { installed: false, supported: true, running: false, busy: false, nodes: 0, subscriptions: 0, phase: 'not_installed', endpoint: 'http://127.0.0.1:3101' }
describe('Mihomo settings', () => {
  beforeEach(() => { vi.resetAllMocks(); get.mockResolvedValue({ data: base }); post.mockResolvedValue({ data: base }) })
  it('installs only after explicit action and does not emit an unready proxy', async () => {
    const wrapper = mount(MihomoSettings); await flushPromises()
    expect(post).not.toHaveBeenCalled()
    await wrapper.findAll('button').find(b => b.text() === '检测并安装')!.trigger('click'); await flushPromises()
    expect(post).toHaveBeenCalledWith('/admin/system/mihomo', expect.objectContaining({ action: 'install' }))
    expect(wrapper.emitted('ready')).toBeUndefined(); wrapper.unmount()
  })
  it('submits multiple subscriptions without putting them in status output', async () => {
    get.mockResolvedValue({ data: { ...base, installed: true } })
    const wrapper = mount(MihomoSettings); await flushPromises()
    await wrapper.get('textarea').setValue('https://example.org/a?token=secret\nhttps://example.org/b')
    await wrapper.findAll('button').find(b => b.text() === '保存并应用')!.trigger('click'); await flushPromises()
    expect(post).toHaveBeenCalledWith('/admin/system/mihomo', expect.objectContaining({ action: 'apply', subscriptions: ['https://example.org/a?token=secret', 'https://example.org/b'] }))
    expect(wrapper.get<HTMLTextAreaElement>('textarea').element.value).toBe(''); wrapper.unmount()
  })
})
