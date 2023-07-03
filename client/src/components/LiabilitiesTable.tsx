import React, { type ReactNode } from 'react'
import {
  Grid,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper
} from '@mui/material'
import axios from 'axios'

interface Liability {
  Name: string
}

export default class LiabilitiesTable extends React.Component {
  state = {
    liabilities: [] as Liability[]
  }

  componentDidMount (): void {
    axios.get('http://localhost:8080/api/liability')
      .then(res => {
        const liabilities: Liability = res.data
        this.setState({ liabilities })
      }).catch(console.error)
  }

  render (): ReactNode {
    return (
    <Grid container spacing = {2} justifyContent="center" alignItems="center">
      <Grid item xs={8}>
        <TableContainer component={Paper}>
          <Table sx={{ minWidth: 650 }} aria-label="simple table">
            <TableHead>
              <TableRow>
                <TableCell><b>Liability Name</b></TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {this.state.liabilities.map((liability) => (
                <TableRow
                  key={liability.Name}
                  sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                >
                  <TableCell component="th" scope="row">{liability.Name}</TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </Grid>
    </Grid>
    )
  }
}
