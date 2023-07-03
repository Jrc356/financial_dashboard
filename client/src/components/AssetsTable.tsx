import React, { type ReactNode } from 'react'
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper
} from '@mui/material'
import axios from 'axios'

interface Asset {
  Name: string
  Type: string
  TaxBucket: string
}

export default class AssetsTable extends React.Component {
  state = {
    assets: [] as Asset[]
  }

  componentDidMount (): void {
    axios.get('http://localhost:8080/api/asset')
      .then(res => {
        const assets: Asset = res.data
        this.setState({ assets })
      }).catch(console.error)
  }

  render (): ReactNode {
    return (
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell><b>Asset Name</b></TableCell>
              <TableCell><b>Type</b></TableCell>
              <TableCell><b>Tax Bucket</b></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {this.state.assets.map((asset) => (
              <TableRow
                key={asset.Name}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row">{asset.Name}</TableCell>
                <TableCell>{asset.Type}</TableCell>
                <TableCell>{asset.TaxBucket}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    )
  }
}
