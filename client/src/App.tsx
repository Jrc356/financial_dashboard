import React from 'react'
import './App.css'
import AssetsTable from './components/desktop/AssetsTable'
import LiabilitiesTable from './components/desktop/LiabilitiesTable'
import {
  Grid
} from '@mui/material'
import { ThemeProvider, createTheme } from '@mui/material/styles'
import CssBaseline from '@mui/material/CssBaseline'
import ValuesTableWithChart from './components/desktop/ValuesTableWithChart'

const darkTheme = createTheme({
  palette: {
    mode: 'dark'
  }
})

export default class App extends React.Component {
  render (): React.ReactNode {
    return (
      <ThemeProvider theme={darkTheme}>
        <CssBaseline />

        <Grid container justifyContent="center" alignItems="center" direction={'column'}>
          <Grid item xs={6} marginTop={3}>
            <ValuesTableWithChart accountType="asset"></ValuesTableWithChart>
          </Grid>

          <Grid item xs={6} marginTop={3}>
            <ValuesTableWithChart accountType="liability"></ValuesTableWithChart>
          </Grid>

          <Grid item xs={6} marginTop={3}>
            <AssetsTable></AssetsTable>
          </Grid>

          <Grid item xs={6} marginTop={3}>
            <LiabilitiesTable></LiabilitiesTable>
          </Grid>
        </Grid>
      </ThemeProvider>
    )
  }
}
