import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'
import HarvestFlowView from '../HarvestFlowView.vue'

const { getFlow } = vi.hoisted(() => ({ getFlow: vi.fn() }))

vi.mock('@/api/admin/accounts', () => ({
  getCodexHarvestFlow: getFlow,
  updateCodexSkipHarvest: vi.fn(),
  updateCodexHarvestConfig: vi.fn()
}))
// vue-i18n 被整体替换，所以模块里用到的导出都必须提供。
// 组件现在会经 @/api/client → @/i18n 间接使用 createI18n，
// 只 mock useI18n 会让 createI18n 变成 undefined，模块初始化即报错。
vi.mock('vue-i18n', () => ({
  useI18n: () => ({ t: (key: string) => key }),
  createI18n: () => ({ global: { t: (key: string) => key, locale: { value: 'zh' } } })
}))
vi.mock('@/stores', () => ({
  useAppStore: () => ({
    showSuccess: vi.fn(),
    showError: vi.fn(),
    showInfo: vi.fn(),
    showWarning: vi.fn()
  })
}))
vi.mock('@/components/layout/AppLayout.vue', () => ({
  default: { template: '<main><slot /></main>' }
}))

function response() {
  return {
    harvest: { enabled: false, fail_closed: false, strategy: 'standby', models: [] },
    sidecar: { reachable: false },
    stages: [],
    accounts: [{
      id: 1,
      name: 'Test account',
      status: 'active',
      schedulable: true,
      ready_count: 0,
      tickets: null
    }],
    counts: {
      tickets_ready: 0, tickets_blocked: 0, probe_hit: 0, probe_miss: 0,
      ticket_accept: 0, ticket_reject: 0, select_ok: 0, select_skip: 0, select_fail: 0
    },
    events: []
  }
}

describe('HarvestFlowView nullable API lists', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    getFlow.mockReset()
  })
  afterEach(() => {
    vi.useRealTimers()
  })

  it.each([
    {
      name: 'disabled harvesting returns null tickets',
      payload: response(),
      text: '0/0'
    },
    {
      name: 'no eligible accounts returns null accounts',
      payload: { ...response(), accounts: null },
      text: 'admin.harvestFlow.noAccounts'
    },
    {
      name: 'null events still render the empty state',
      payload: { ...response(), accounts: [], events: null },
      text: 'admin.harvestFlow.noEvents'
    },
    {
      name: 'normal ticket lists still render',
      payload: {
        ...response(),
        accounts: [{
          ...response().accounts[0],
          tickets: [{ model: 'gpt-5.6-sol', ready: false, blocked: false, remaining_seconds: 0 }]
        }]
      },
      text: 'gpt-5.6-sol'
    }
  ])('$name', async ({ payload, text }) => {
    getFlow.mockResolvedValue(payload)
    const errors: unknown[] = []
    const wrapper = mount(HarvestFlowView, {
      global: {
        stubs: { Icon: true, LoadingSpinner: true },
        config: { errorHandler: (error) => { errors.push(error) } }
      }
    })
    try {
      await flushPromises()
      expect(errors).toEqual([])
      expect(wrapper.text()).toContain('admin.harvestFlow.title')
      expect(wrapper.text()).toContain(text)
    } finally {
      wrapper.unmount()
    }
  })
})
