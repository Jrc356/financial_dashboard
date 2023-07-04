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

interface Liability {
  Name: string
}

export default function LiabilitiesTable (): JSX.Element {
  const [liabilities, setLiabilities] = React.useState([] as Liability[])

  React.useEffect(() => {
    axios.get('http://localhost:8080/api/liability')
      .then((response) => {
        setLiabilities(response.data)
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
              <TableRow key={asset.Name} sx={{ '&:last-child td, &:last-child th': { border: 0 } }}>
                <TableCell component="th" scope="row">{asset.Name}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    </div>
  )
}