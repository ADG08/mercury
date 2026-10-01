import http from 'k6/http'
import { check, sleep } from 'k6'

export const options = {
    stages: [
        { duration: '10s', target: 50},
        { duration: '30s', target: 150},
        { duration: '10s', target: 0}
    ],
    thresholds: {
        'http_req_failed{phase:read}': ['rate<0.01'],
        'http_req_duration{phase:read}': ['p(95)<10', 'p(99)<25'],
        'checks{phase:read}': ['rate>0.99']
    }
}

export function setup() {
    const params = {
        headers: { 'Content-Type': 'application/json' },
        tags: { phase: 'setup' }
    }
    const targetUrls = [
        'https://www.linkedin.com/in/andre-domingues-ramos/',
        'https://github.com/ADG08',
        'https://www.google.com/',
        'https://stackoverflow.com/',
        'https://andredr.fr/'
    ]
    const readWeights = [50, 25, 15, 7, 3]
    const codes = []

    for (const targetUrl of targetUrls) {
        const payload = JSON.stringify({
            url: targetUrl
        })
        const res = http.post("http://host.docker.internal:8080/url", payload, params)

        if (res.status !== 201) {
            throw new Error(`Echec setup POST /url: code ${res.status} - body: ${res.body}`)
        }

        const data = JSON.parse(res.body)
        const parts = data.short_url.split('/')
        codes.push(parts[parts.length - 1])
    }

    return { codes: codes, readWeights: readWeights }
}

export default function (data) {
    const randomValue = Math.random() * 100
    let accumulatedWeight = 0
    let code = data.codes[data.codes.length - 1]

    for (let index = 0; index < data.codes.length; index++) {
        accumulatedWeight += data.readWeights[index]
        if (randomValue < accumulatedWeight) {
            code = data.codes[index]
            break
        }
    }

    const res = http.get(`http://host.docker.internal:8080/${code}`, {
        redirects: 0,
        tags: { phase: 'read' }
    })

    check(res, {
        'statut 307': (r) => r.status === 307
    }, { phase: 'read' })

    sleep(0.05)
}