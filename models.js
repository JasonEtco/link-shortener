import { URL } from 'url'
import { LinksTable as SQLinksTable, StatsTable as SQStatsTable } from './models.sqlite'
import { LinksTable as MemLinksTable, StatsTable as MemStatsTable } from './models.memory'
import createDebug from 'debug'
import dotenv from 'dotenv'

const debug = createDebug('models')
dotenv.config()

export class LinksTable {
  constructor(url) {
    url = new URL(url)
    if (url.protocol === 'sqlite:') return new SQLinksTable(url.toString())
    if (url.protocol === 'memory:') return new MemLinksTable(url.toString())
    throw new Error(`unknown database protocol ${url.protocol}`)
  }
}

export class StatsTable {
  constructor(url) {
    url = new URL(url)
    if (url.protocol === 'sqlite:') return new SQStatsTable(url.toString())
    if (url.protocol === 'memory:') return new MemStatsTable(url.toString())
    throw new Error(`unknown database protocol ${url.protocol}`)
  }
}

if (require.main === module) {
  debug(`Migrating ${process.env.DB}`)
  Promise.all([new LinksTable(process.env.DB).migrate(), new StatsTable(process.env.DB).migrate()]).then(() =>
    console.log(`Migrations complete`)
  )
}
