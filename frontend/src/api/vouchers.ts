export type Voucher = {
  crewName: string
  crewId: string
  flightNumber: string
  flightDate: string
  aircraft: string
  seats: string[]
}

export type CheckVoucherInput = {
  flightNumber: string
  flightDate: string
}

export type GenerateVoucherInput = CheckVoucherInput & {
  crewName: string
  crewId: string
  aircraft: string
}

type CheckVoucherResponse = {
  exists: boolean
  voucher: Voucher | null
}

type ApiErrorResponse = {
  error?: {
    message?: string
  }
}

async function postJSON<T>(url: string, body: object): Promise<T> {
  let response: Response

  try {
    response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
  } catch {
    throw new Error('Unable to reach the voucher service. Please try again.')
  }

  if (!response.ok) {
    let message = `Voucher service returned status ${response.status}.`

    try {
      const payload = (await response.json()) as ApiErrorResponse
      if (payload.error?.message) {
        message = payload.error.message
      }
    } catch {
      // Keep the status-based message when the response is not JSON.
    }

    throw new Error(message)
  }

  try {
    return (await response.json()) as T
  } catch {
    throw new Error('Voucher service returned an invalid response.')
  }
}

export function checkVoucher(input: CheckVoucherInput) {
  return postJSON<CheckVoucherResponse>('/api/check', input)
}

export function generateVoucher(input: GenerateVoucherInput) {
  return postJSON<Voucher>('/api/generate', input)
}
