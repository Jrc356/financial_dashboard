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
import { GetAllAccountsByClass, type Account } from '../lib/api'

export default function LiabilitiesTable (): JSX.Element {
  const [liabilities, setLiabilities] = React.useState([] as Account[])

  React.useEffect(() => {
    GetAllAccountsByClass('liability')
      .then((l) => {
        setLiabilities(l)
      })
      .catch(console.error)
  }, [])

  if (liabilities.length === 0) return <div></div>

  return (
    <div>
      <TableContainer component={Paper}>
        <Table sx={{ minWidth: 650 }} aria-label="simple table">
          <TableHead>
            <TableRow>
              <TableCell><b>Liability Name</b></TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {liabilities.map((asset) => (
              <TableRow key={asset.name} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                <TableCell component="th" scope="row">{asset.name}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  )
}
