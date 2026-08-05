import { Container, Stack, Text, Title } from '@mantine/core'

function App() {
  return (
    <Container component="main" size="sm" py="xl">
      <Stack gap="md">
        <Title order={1}>Airline Voucher Seat Assignment</Title>
        <Text c="dimmed">
          A simple foundation for managing airline voucher seat assignments.
        </Text>
      </Stack>
    </Container>
  )
}

export default App
