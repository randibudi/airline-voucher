import { MantineProvider } from '@mantine/core'
import { render as testingLibraryRender } from '@testing-library/react'
import type { ReactElement } from 'react'

function render(ui: ReactElement) {
  return testingLibraryRender(ui, {
    wrapper: ({ children }) => (
      <MantineProvider env="test">{children}</MantineProvider>
    ),
  })
}

export * from '@testing-library/react'
export { render }
