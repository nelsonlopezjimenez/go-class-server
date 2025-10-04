import { Alert as ChakraAlert } from '@chakra-ui/react'
import * as React from 'react'

export const Alert = React.forwardRef(function Alert(props, ref) {
  const { title, children, icon, startElement, endElement, ...rest } = props
  return (
    <ChakraAlert.Root ref={ref} {...rest}>
      {startElement || <ChakraAlert.Indicator>{icon}</ChakraAlert.Indicator>}
      {children ? (
        <ChakraAlert.Content>
          <ChakraAlert.Title>{title}</ChakraAlert.Title>
          <ChakraAlert.Description>{children}</ChakraAlert.Description>
        </ChakraAlert.Content>
      ) : (
        <ChakraAlert.Title flex='1'>{title}</ChakraAlert.Title>
      )}
      {endElement}
    </ChakraAlert.Root>
  )
})
