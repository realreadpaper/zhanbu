import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { afterEach, beforeEach, expect, test } from 'vitest'
import Login from './Login'
import { useAuthStore } from '../stores/authStore'

const resetAuthStore = () => {
  useAuthStore.setState({
    user: null,
    isAuthenticated: false,
    isLoading: false,
    error: null,
    errorCode: null,
  })
}

beforeEach(() => {
  localStorage.clear()
  resetAuthStore()
})

afterEach(() => {
  localStorage.clear()
  resetAuthStore()
})

test('clears stale email verification errors when the login page opens', async () => {
  useAuthStore.setState({
    error: '请先验证邮箱后再登录',
    errorCode: 2008,
  })

  render(
    <MemoryRouter>
      <Login />
    </MemoryRouter>
  )

  await waitFor(() => {
    expect(screen.queryByText('请先验证邮箱后再登录')).toBeNull()
  })
  expect(useAuthStore.getState().error).toBeNull()
  expect(useAuthStore.getState().errorCode).toBeNull()
})
