import CssBaseline from '@mui/material/CssBaseline'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import React from 'react'
import { green } from '@mui/material/colors'
import AccountView from './views/AccountView'
import { API } from './lib/api'
import { Box } from '@mui/material'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

export default class App extends React.Component {
  render (): React.ReactNode {
    return (
      <Box style={{
        backgroundColor: green[300],
        textAlign: 'center'
      }}>
        <ThemeProvider theme={darkTheme}>
          <CssBaseline />
          <AccountView accountType={API.Asset}></AccountView>
        </ThemeProvider>
      </Box>
    )
  }
}
