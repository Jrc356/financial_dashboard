import React from 'react'
import Paper from '@mui/material/Paper'
import BottomNavigation from '@mui/material/BottomNavigation'
import BottomNavigationAction from '@mui/material/BottomNavigationAction'
import PaidIcon from '@mui/icons-material/Paid'
import CreditCardIcon from '@mui/icons-material/CreditCard'
import { useLocation, Link } from 'react-router-dom'

export default function BottomNav (): React.ReactElement {
  const pathname = useLocation().pathname
  const [value, setValue] = React.useState(pathname)

  // TODO: fucking router bs

  return (
    <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} elevation={3}>
      <BottomNavigation
        showLabels
        value={value}
        onChange={(event, newValue): void => {
          setValue(newValue)
        }}
      >
        <BottomNavigationAction
          label='Assets'
          icon={<PaidIcon/>}
          component={Link}
          to='/assets'
          value='/assets'
        />
        <BottomNavigationAction
          label='Liabilities'
          icon={<CreditCardIcon/>}
          component={Link}
          to='/liabilities'
          value='/liabilities'
        />
      </BottomNavigation>
    </Paper>
  )
}
