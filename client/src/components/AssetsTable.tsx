import {
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow
} from '@mui/material'
import React from 'react'
import { type Account, GetAllAccountsByClass } from '../lib/api'

export default function AssetsTable (): JSX.Element {
  const [assets, setAssets] = React.useState([] as Account[])

  React.useEffect(() => {
    GetAllAccountsByClass('asset')
      .then((a) => {
        setAssets(a)
      })
      .catch(console.error)
  }, [])

  if (assets.length === 0) return <div><TableContainer component={Paper}></TableContainer></div>

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
                key={asset.name}
                sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
              >
                <TableCell component="th" scope="row">{asset.name}</TableCell>
                <TableCell>{asset.class}</TableCell>
                <TableCell>{asset.taxBucket}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  )
}
