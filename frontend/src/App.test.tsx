import { screen } from './test/render'
import { render } from './test/render'
import { describe, expect, it } from 'vitest'
import App from './App'

describe('App', () => {
  it('renders the application heading', () => {
    render(<App />)

    expect(
      screen.getByRole('heading', {
        level: 1,
        name: 'Airline Voucher Seat Assignment',
      }),
    ).toBeInTheDocument()
  })
})
