import { QrCode as ChakraQrCode } from '@chakra-ui/react'
import * as React from 'react'

export const QrCode = React.forwardRef(function QrCode(props, ref) {
  const { children, fill, overlay, ...rest } = props
  return (
    <ChakraQrCode.Root ref={ref} {...rest}>
      <ChakraQrCode.Frame style={{ fill }}>
        <ChakraQrCode.Pattern />
      </ChakraQrCode.Frame>
      {overlay}
      {children && <ChakraQrCode.Overlay>{children}</ChakraQrCode.Overlay>}
    </ChakraQrCode.Root>
  )
})
