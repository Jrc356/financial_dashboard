import axios from 'axios'
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper
} from '@mui/material'
import React from 'react'

export interface Asset {
  Name: string
  Type: string
  TaxBucket: string
}

export default function AssetsTable (): JSX.Element {
  const [assets, setAssets] = React.useState([] as Asset[])

  React.useEffect(() => {
    axios.get('http://localhost:8080/api/asset')
      .then((response) => {
        setAssets(response.data)
      })
      .catch(console.error)
  }, [])

  if (assets.length === 0) return <div></div>

  return (
    <div>
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
            {assets.map((asset) => (
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
    </div>
  )
}
