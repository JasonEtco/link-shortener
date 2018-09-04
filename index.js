const express = require('express')
const sqlite = require('sqlite')
const bodyParser = require('body-parser')

const dbPromise = sqlite.open('./urls.db', { Promise })
const port = process.env.PORT || 3000
const app = express()

if (!process.env.AUTH_HEADER) {
  console.log('AUTH_HEADER not set')
  process.exit(1)
}

app.post('/', bodyParser.text(), async (req, res) => {
  if (req.header('auth') !== process.env.AUTH_HEADER) {
    return res.status(401).send('Authentication failed!')
  }
  if (typeof req.body !== 'string') {
    return res.status(400).send('Expecting Content-Type: text/plain')
  }
  if (!req.body.includes('http')) {
    return res.status(400).send('Protocol missing')
  }
  const db = await dbPromise
  let results = await db.get('SELECT id FROM urls WHERE url = ?', req.body)
  if (results) {
    return res.redirect(`${results.id}`)
  }

  let id = Math.random().toString(36).slice(2)
  results = await db.get('SELECT id FROM urls WHERE id = ?', id)
  while(results) {
    results = await db.get('SELECT id FROM urls WHERE id = ?', id)
    id = Math.random().toString(36).slice(2)
  }

  await db.run('INSERT INTO urls VALUES(?, ?);', id, req.body)
  res.redirect(`${id}`)
})

app.post('/:id', bodyParser.text(), async (req, res) => {
  if (req.header('auth') !== process.env.AUTH_HEADER) {
    return res.status(401).send('Authentication failed!')
  }
  if (typeof req.body !== 'string') {
    return res.status(400).send('Expecting Content-Type: text/plain')
  }
  if (!req.body.includes('http')) {
    return res.status(400).send('Protocol missing')
  }
  const db = await dbPromise
  const url = await db.get('SELECT url FROM urls WHERE id = ?', req.params.id)
  if (url) {
    res.status(409).send('Already a URL with that id!')
  } else {
    await db.run('INSERT INTO urls VALUES(?, ?);', req.params.id, req.body)
    res.redirect(`${req.params.id}`)
  }
})

app.get('/:id', async (req, res) => {
  const db = await dbPromise
  const results = await db.get('SELECT url FROM urls WHERE id = ?', req.params.id)
  if (results) {
    res.redirect(results.url)
    await db.run('INSERT INTO stats VALUES(?, ?, ?);', req.params.id, 200, JSON.stringify(req.headers))
  } else {
    res.status(404).send("Sorry can't find that!")
    await db.run('INSERT INTO stats VALUES(?, ?, ?);', req.params.id, 400, JSON.stringify(req.headers))
  }
})

app.listen(port)
