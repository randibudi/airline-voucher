import {
  Alert,
  Badge,
  Button,
  Container,
  Divider,
  Grid,
  Loader,
  Paper,
  Select,
  SimpleGrid,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core'
import { DateInput } from '@mantine/dates'
import { useForm } from '@mantine/form'
import dayjs from 'dayjs'
import customParseFormat from 'dayjs/plugin/customParseFormat'
import { useState } from 'react'
import {
  checkVoucher,
  generateVoucher,
  type GenerateVoucherInput,
  type Voucher,
} from './api/vouchers'
import './index.css'

dayjs.extend(customParseFormat)

const aircraftOptions = ['ATR', 'Airbus 320', 'Boeing 737 Max']
const requiredMessage = 'This field is required'

function parseFlightDate(value: string) {
  const parsed = dayjs(value, 'DD-MM-YYYY', true)
  return parsed.isValid() ? parsed.toDate() : null
}

type FormValues = {
  crewName: string
  crewId: string
  flightNumber: string
  flightDate: Date | null
  aircraft: string | null
}

type ResultState = {
  voucher: Voucher
  existing: boolean
}

function VoucherResult({ voucher, existing }: ResultState) {
  return (
    <Paper
      component="section"
      aria-live="polite"
      aria-atomic="true"
      withBorder
      radius="lg"
      p={{ base: 'lg', sm: 'xl' }}
      className="result-card"
    >
      <Stack gap="lg">
        <div>
          <Badge
            color={existing ? 'blue' : 'teal'}
            variant="light"
            size="lg"
            mb="sm"
          >
            {existing ? 'Existing voucher' : 'New voucher generated'}
          </Badge>
          <Title order={2} size="h3">
            Voucher details
          </Title>
          <Text c="dimmed" mt={4} size="sm">
            {existing
              ? 'A voucher was already created for this flight and date.'
              : 'Three seats have been assigned and saved successfully.'}
          </Text>
        </div>

        <Divider />

        <SimpleGrid
          cols={{ base: 1, xs: 2 }}
          spacing={{ base: 'md', sm: 'lg' }}
        >
          <VoucherField label="Crew Name" value={voucher.crewName} />
          <VoucherField label="Crew ID" value={voucher.crewId} />
          <VoucherField label="Flight Number" value={voucher.flightNumber} />
          <VoucherField
            label="Flight Date"
            value={dayjs(voucher.flightDate).format('DD-MM-YYYY')}
          />
          <VoucherField label="Aircraft" value={voucher.aircraft} />
        </SimpleGrid>

        <Divider />

        <div>
          <Title order={3} size="h4" mb="sm">
            Assigned Seats
          </Title>
          <div className="seat-grid" role="list" aria-label="Assigned seats">
            {voucher.seats.map((seat) => (
              <Paper
                key={seat}
                role="listitem"
                aria-label={`Assigned seat ${seat}`}
                withBorder
                radius="md"
                className="seat-card"
              >
                <Text className="seat-label">Seat</Text>
                <Text className="seat-value">{seat}</Text>
              </Paper>
            ))}
          </div>
        </div>
      </Stack>
    </Paper>
  )
}

function VoucherField({ label, value }: { label: string; value: string }) {
  return (
    <div className="voucher-field">
      <Text className="field-label">{label}</Text>
      <Text fw={650} className="field-value">
        {value}
      </Text>
    </div>
  )
}

function HowItWorks() {
  const steps = [
    'We first check for a voucher with the same flight number and date.',
    'If none exists, three available seats are assigned automatically.',
    'The voucher is saved for future checks of that flight and date.',
  ]

  return (
    <Paper
      component="aside"
      withBorder
      radius="lg"
      p={{ base: 'lg', sm: 'xl' }}
      className="info-card"
    >
      <Badge variant="light" color="blue" mb="sm">
        Quick guide
      </Badge>
      <Title order={2} size="h3">
        How it works
      </Title>
      <Text c="dimmed" size="sm" mt={6} mb="lg">
        One request is all it takes to find or create your voucher.
      </Text>
      <Stack component="ol" gap="md" className="steps-list">
        {steps.map((step, index) => (
          <li key={step} className="step-item">
            <span className="step-number" aria-hidden="true">
              {index + 1}
            </span>
            <Text size="sm">{step}</Text>
          </li>
        ))}
      </Stack>
    </Paper>
  )
}

function getErrorMessage(caught: unknown) {
  if (
    caught instanceof Error &&
    caught.message === 'Unable to reach the voucher service. Please try again.'
  ) {
    return caught.message
  }

  return 'The voucher service could not complete the request. Please check your details and try again.'
}

function App() {
  const [result, setResult] = useState<ResultState | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const form = useForm<FormValues>({
    mode: 'uncontrolled',
    initialValues: {
      crewName: '',
      crewId: '',
      flightNumber: '',
      flightDate: null,
      aircraft: null,
    },
    validate: {
      crewName: (value) => (value.trim() ? null : requiredMessage),
      crewId: (value) => (value.trim() ? null : requiredMessage),
      flightNumber: (value) => (value.trim() ? null : requiredMessage),
      flightDate: (value) => (value ? null : requiredMessage),
      aircraft: (value) => (value ? null : requiredMessage),
    },
  })

  const submit = async (values: FormValues) => {
    setResult(null)
    setError(null)
    setLoading(true)

    const request: GenerateVoucherInput = {
      crewName: values.crewName.trim(),
      crewId: values.crewId.trim(),
      flightNumber: values.flightNumber.trim().toUpperCase(),
      flightDate: dayjs(values.flightDate).format('YYYY-MM-DD'),
      aircraft: values.aircraft!,
    }

    try {
      const checked = await checkVoucher({
        flightNumber: request.flightNumber,
        flightDate: request.flightDate,
      })

      if (checked.exists) {
        if (!checked.voucher) {
          throw new Error('Voucher service returned an invalid response.')
        }
        setResult({ voucher: checked.voucher, existing: true })
        return
      }

      const voucher = await generateVoucher(request)
      setResult({ voucher, existing: false })
    } catch (caught) {
      setError(getErrorMessage(caught))
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app-shell">
      <Container component="main" size="lg" py={{ base: 32, sm: 48, md: 64 }}>
        <header className="hero">
          <Text className="eyebrow">CREW VOUCHER</Text>
          <Title order={1} className="hero-title">
            Airline Voucher Seat Assignment
          </Title>
          <Text className="hero-description">
            Check an existing voucher or generate three available seats for your
            flight.
          </Text>
        </header>

        <Grid>
          <Grid.Col span={{ base: 12, md: 7 }}>
            <Paper
              component="section"
              withBorder
              radius="lg"
              p={{ base: 'lg', sm: 32 }}
              className="form-card"
            >
              <Title order={2} size="h3">
                Flight and Crew Details
              </Title>
              <Text c="dimmed" size="sm" mt={6} mb="xl">
                Enter the crew and flight information below. All fields are
                required.
              </Text>

              <form onSubmit={form.onSubmit(submit)} noValidate>
                <Stack gap="lg">
                  <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="lg">
                    <TextInput
                      label="Crew Name"
                      placeholder="e.g. John Doe"
                      description="Enter the name shown on the crew record."
                      required
                      size="md"
                      autoComplete="name"
                      key={form.key('crewName')}
                      {...form.getInputProps('crewName')}
                    />
                    <TextInput
                      label="Crew ID"
                      placeholder="e.g. CRW001"
                      description="Use your assigned crew identifier."
                      required
                      size="md"
                      autoComplete="username"
                      key={form.key('crewId')}
                      {...form.getInputProps('crewId')}
                    />
                    <TextInput
                      label="Flight Number"
                      placeholder="e.g. GA123"
                      description="Spaces are trimmed and letters are capitalized."
                      required
                      size="md"
                      autoComplete="off"
                      autoCapitalize="characters"
                      key={form.key('flightNumber')}
                      {...form.getInputProps('flightNumber')}
                    />
                    <DateInput
                      label="Flight Date"
                      placeholder="DD-MM-YYYY"
                      description="Use day-month-year format."
                      valueFormat="DD-MM-YYYY"
                      dateParser={parseFlightDate}
                      required
                      clearable
                      size="md"
                      autoComplete="off"
                      key={form.key('flightDate')}
                      {...form.getInputProps('flightDate')}
                    />
                  </SimpleGrid>

                  <Select
                    label="Aircraft"
                    placeholder="Select an aircraft"
                    description="Choose the aircraft scheduled for this flight."
                    data={aircraftOptions}
                    required
                    size="md"
                    key={form.key('aircraft')}
                    {...form.getInputProps('aircraft')}
                  />

                  {error && (
                    <Alert
                      color="red"
                      title="Unable to process request"
                      role="alert"
                      variant="light"
                    >
                      {error}
                    </Alert>
                  )}

                  <div className="request-status" aria-live="polite">
                    {loading && (
                      <div className="loading-message">
                        <Loader size="sm" aria-hidden="true" />
                        <Text size="sm" fw={600}>
                          Checking voucher availability...
                        </Text>
                      </div>
                    )}
                  </div>

                  <Button
                    type="submit"
                    loading={loading}
                    disabled={loading}
                    size="md"
                    className="submit-button"
                  >
                    Check or Generate Voucher
                  </Button>
                </Stack>
              </form>
            </Paper>
          </Grid.Col>

          <Grid.Col span={{ base: 12, md: 5 }}>
            <Stack gap="lg">
              <HowItWorks />
              {result && <VoucherResult {...result} />}
            </Stack>
          </Grid.Col>
        </Grid>
      </Container>
    </div>
  )
}

export default App
