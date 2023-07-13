import React from 'react'
import { API, GetAll, type Account, GetValuesForAccount } from '../lib/api'
import {
  Divider,
  Grid,
  Typography
} from '@mui/material'
import AccountCard from '../components/AccountCard'
import moneyFormatter from '../lib/formatter'

function toTitleCase (str: string): string {
  return str.charAt(0).toUpperCase() + str.substr(1).toLowerCase()
}

interface Props {
  accountType: API
}

export default function AccountView ({
  accountType
}: Props): React.ReactElement {
  const [accounts, setAccounts] = React.useState([] as Account[])
  const [totalValue, setTotalValue] = React.useState(0)

  let api = API.Asset
  switch (accountType) {
    case API.Asset:
      api = API.Asset
      break
    case API.Liability:
      api = API.Liability
      break
    default:
      console.error('unknown account type')
      break
  }

  // TODO: This is so gross
  React.useEffect(() => {
    GetAll(api)
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
  }, [])

  React.useEffect(() => {
    if (accounts.length === 0) {
      return
    }
    const valuePromises = []
    for (const account of accounts) {
      if (account.Values.length === 0) {
        valuePromises.push(GetValuesForAccount(api, account))
      }
    }
    if (valuePromises.length > 0) {
      Promise.all(valuePromises)
        .then((accs) => {
          setAccounts(accs)
        })
        .catch(console.error)
    }
  }, [accounts])

  React.useEffect(() => {
    if (accounts.length === 0) {
      return
    }
    let v = 0
    for (const account of accounts) {
      if (account.Values.length > 0) {
        v += account.Values[0].Value
      }
    }
    setTotalValue(v)
  }, [accounts])

  return (
    <Grid
      container
      justifyContent="center"
      alignItems="center"
      direction={'column'}
      overflow={'auto'}
      flex={1}
    >
      <Grid item>
        <Typography variant="h6" color={'black'} marginTop={6}>
          Total {toTitleCase(api)} Value:
        </Typography>
        <Typography variant="h2" color={'black'}>{moneyFormatter.format(totalValue)}</Typography>
        <Divider
          sx={{
            marginTop: 6,
            marginBottom: 6,
            border: 2,
            borderColor: 'black',
            opacity: 100,
            marginLeft: '5%',
            marginRight: '5%',
            justifyContent: 'center'
          }}
        >
        </Divider>
        {accounts.map((account, i) => (
          <AccountCard key={i} account={account}></AccountCard>
        ))}
      </Grid>
    </Grid>
  )
}
