import {
  Alert,
  Badge,
  Button,
  Container,
  Grid,
  Group,
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
    <Paper component="section" aria-live="polite" withBorder radius="md" p="lg">
      <Stack gap="lg">
        <Alert color={existing ? 'blue' : 'green'} title="Voucher status">
          {existing
            ? 'A previously issued voucher was found.'
            : 'A new voucher was created successfully.'}
        </Alert>

        <div>
          <Title order={2} size="h3">
            Voucher details
          </Title>
          <Text c="dimmed" size="sm">
            Seat assignment for {voucher.flightNumber}
          </Text>
        </div>

        <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="md">
          <VoucherField label="Crew Name" value={voucher.crewName} />
          <VoucherField label="Crew ID" value={voucher.crewId} />
          <VoucherField label="Flight Number" value={voucher.flightNumber} />
          <VoucherField
            label="Flight Date"
            value={dayjs(voucher.flightDate).format('DD-MM-YYYY')}
          />
          <VoucherField label="Aircraft" value={voucher.aircraft} />
        </SimpleGrid>

        <div>
          <Text fw={600} mb="xs">
            Assigned Seats
          </Text>
          <Group gap="sm" aria-label="Assigned seats">
            {voucher.seats.map((seat) => (
              <Badge key={seat} size="xl" variant="filled" radius="sm">
                {seat}
              </Badge>
            ))}
          </Group>
        </div>
      </Stack>
    </Paper>
  )
}

function VoucherField({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <Text size="sm" c="dimmed">
        {label}
      </Text>
      <Text fw={600}>{value}</Text>
    </div>
  )
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
          throw new Error('Voucher service returned an invalid check response.')
        }
        setResult({ voucher: checked.voucher, existing: true })
        return
      }

      const voucher = await generateVoucher(request)
      setResult({ voucher, existing: false })
    } catch (caught) {
      setError(
        caught instanceof Error
          ? caught.message
          : 'An unexpected error occurred. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <Container component="main" size="md" py={{ base: 'lg', sm: 'xl' }}>
      <Stack gap="xl">
        <div>
          <Title order={1}>Airline Voucher Seat Assignment</Title>
          <Text c="dimmed" mt="xs">
            Check for an existing voucher or generate a new seat assignment.
          </Text>
        </div>

        <Paper withBorder radius="md" p={{ base: 'md', sm: 'xl' }}>
          <form onSubmit={form.onSubmit(submit)} noValidate>
            <Stack gap="md">
              <Grid>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <TextInput
                    label="Crew Name"
                    placeholder="Jane Doe"
                    required
                    autoComplete="name"
                    key={form.key('crewName')}
                    {...form.getInputProps('crewName')}
                  />
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <TextInput
                    label="Crew ID"
                    placeholder="CRW001"
                    required
                    key={form.key('crewId')}
                    {...form.getInputProps('crewId')}
                  />
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <TextInput
                    label="Flight Number"
                    placeholder="GA123"
                    required
                    autoCapitalize="characters"
                    key={form.key('flightNumber')}
                    {...form.getInputProps('flightNumber')}
                  />
                </Grid.Col>
                <Grid.Col span={{ base: 12, sm: 6 }}>
                  <DateInput
                    label="Flight Date"
                    placeholder="DD-MM-YYYY"
                    valueFormat="DD-MM-YYYY"
                    dateParser={parseFlightDate}
                    required
                    clearable
                    key={form.key('flightDate')}
                    {...form.getInputProps('flightDate')}
                  />
                </Grid.Col>
                <Grid.Col span={12}>
                  <Select
                    label="Aircraft"
                    placeholder="Select an aircraft"
                    data={aircraftOptions}
                    required
                    key={form.key('aircraft')}
                    {...form.getInputProps('aircraft')}
                  />
                </Grid.Col>
              </Grid>

              {error && (
                <Alert
                  color="red"
                  title="Unable to process voucher"
                  role="alert"
                >
                  {error}
                </Alert>
              )}

              <Button type="submit" loading={loading} disabled={loading}>
                Check and Generate Voucher
              </Button>
            </Stack>
          </form>
        </Paper>

        {result && <VoucherResult {...result} />}
      </Stack>
    </Container>
  )
}

export default App
