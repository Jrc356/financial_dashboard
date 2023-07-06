import React from 'react'
import { API, GetAll, type Account, GetValuesForAccount } from '../../lib/api'
import {
  Card,
  CardActionArea,
  CardContent,
  Divider,
  Grid,
  Typography
} from '@mui/material'
import { green } from '@mui/material/colors'

const formatter = new Intl.NumberFormat('en-US', {
  style: 'currency',
  currency: 'USD'
})

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
      } else {
        valuePromises.push(
          new Promise<Account>((resolve, reject) => {
            resolve(account)
          })
        )
      }
    }
    Promise.all(valuePromises)
      .then((accs) => {
        setAccounts(accs)
      })
      .catch(console.error)
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
      overflow={'hidden'}
    >
      <Grid item>
        <Typography variant="h6" marginTop={6}>
          Total Value:
        </Typography>
        <Typography variant="h2">{formatter.format(totalValue)}</Typography>
        <Divider
          sx={{
            marginTop: 6,
            marginBottom: 6,
            border: 2,
            borderColor: 'white',
            opacity: 100,
            marginLeft: '5%',
            marginRight: '5%',
            justifyContent: 'center'
          }}
        ></Divider>
        {accounts.map((account, i) => (
          <Card
            key={i}
            sx={{
              backgroundColor: 'white',
              height: 120,
              width: 330,
              margin: 2,
              borderRadius: 3
            }}
          >
            <CardActionArea onClick={() => { console.log(account.Name) }} TouchRippleProps={{ color: green[100] }}>
              <Card
                sx={{
                  backgroundColor: 'black',
                  height: 60,
                  width: 165,
                  margin: 2,
                  borderRadius: 3
                }}
              ></Card>
              {/* //TODO why the fuck won't the ripple work right */}
              <CardContent>
                {/* <Typography
                  variant="h4"
                  color={'black'}
                  align="left"
                  fontSize={28}
                >
                  {account.Name.length > 18
                    ? `${account.Name.slice(0, 18)}...`
                    : account.Name}
                </Typography>
                <Typography
                  variant="h4"
                  color={'black'}
                  align="left"
                  fontSize={20}
                  paddingTop={3}
                >
                  {account.Values.length > 0 &&
                    formatter.format(account.Values[0].Value)}
                </Typography> */}
              </CardContent>
            </CardActionArea>
          </Card>
        ))}
      </Grid>
    </Grid>
  )
}
