import userEvent from '@testing-library/user-event'
import { cleanup, fireEvent, render, screen, waitFor } from './test/render'
import { afterEach, describe, expect, it, vi } from 'vitest'
import App from './App'

const voucher = {
  crewName: 'Jane Doe',
  crewId: 'CRW001',
  flightNumber: 'GA123',
  flightDate: '2026-08-05',
  aircraft: 'ATR',
  seats: ['1A', '8C', '18F'],
}

function jsonResponse(body: object) {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { 'Content-Type': 'application/json' },
  })
}

async function completeForm() {
  const user = userEvent.setup()

  await user.type(
    screen.getByRole('textbox', { name: 'Crew Name' }),
    'Jane Doe',
  )
  await user.type(screen.getByRole('textbox', { name: 'Crew ID' }), 'CRW001')
  await user.type(
    screen.getByRole('textbox', { name: 'Flight Number' }),
    ' ga123 ',
  )
  fireEvent.change(screen.getByRole('textbox', { name: 'Flight Date' }), {
    target: { value: '05-08-2026' },
  })
  await user.click(screen.getByRole('combobox', { name: 'Aircraft' }))
  await user.click(screen.getByRole('option', { name: 'ATR' }))
  await user.click(
    screen.getByRole('button', { name: 'Check or Generate Voucher' }),
  )
}

afterEach(() => {
  cleanup()
  vi.restoreAllMocks()
})

describe('App', () => {
  it('renders the main voucher form fields', () => {
    render(<App />)

    expect(
      screen.getByRole('textbox', { name: 'Crew Name' }),
    ).toBeInTheDocument()
    expect(screen.getByRole('textbox', { name: 'Crew ID' })).toBeInTheDocument()
    expect(
      screen.getByRole('textbox', { name: 'Flight Number' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('textbox', { name: 'Flight Date' }),
    ).toBeInTheDocument()
    expect(
      screen.getByRole('combobox', { name: 'Aircraft' }),
    ).toBeInTheDocument()
  })

  it('displays an existing voucher without generating a new one', async () => {
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(jsonResponse({ exists: true, voucher }))
    render(<App />)

    await completeForm()

    expect(await screen.findByText('Existing voucher')).toBeInTheDocument()
    expect(
      screen.getByText(
        'A voucher was already created for this flight and date.',
      ),
    ).toBeInTheDocument()
    expect(screen.getByText('05-08-2026')).toBeInTheDocument()
    expect(
      screen.getByRole('heading', { name: 'Assigned Seats' }),
    ).toBeInTheDocument()
    expect(screen.getByText('1A')).toBeInTheDocument()
    expect(fetchMock).toHaveBeenCalledTimes(1)
    expect(fetchMock).toHaveBeenCalledWith(
      '/api/check',
      expect.objectContaining({
        body: JSON.stringify({
          flightNumber: 'GA123',
          flightDate: '2026-08-05',
        }),
      }),
    )
  })

  it('generates and displays a new voucher after an empty check', async () => {
    const fetchMock = vi
      .spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(jsonResponse({ exists: false, voucher: null }))
      .mockResolvedValueOnce(jsonResponse(voucher))
    render(<App />)

    await completeForm()

    expect(await screen.findByText('New voucher generated')).toBeInTheDocument()
    expect(
      screen.getByText(
        'Three seats have been assigned and saved successfully.',
      ),
    ).toBeInTheDocument()
    expect(screen.getByText('8C')).toBeInTheDocument()
    await waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(2))
    expect(fetchMock.mock.calls[1]?.[0]).toBe('/api/generate')
    expect(fetchMock.mock.calls[1]?.[1]).toEqual(
      expect.objectContaining({
        body: JSON.stringify({
          crewName: 'Jane Doe',
          crewId: 'CRW001',
          flightNumber: 'GA123',
          flightDate: '2026-08-05',
          aircraft: 'ATR',
        }),
      }),
    )
  })
})
