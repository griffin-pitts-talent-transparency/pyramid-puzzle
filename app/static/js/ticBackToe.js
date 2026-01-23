// ticBackToe.js
import { getUserAsync } from './auth.js';

let MOVE_QUEUE = []
let WINNER = false

function resetGame() {
    WINNER = false
    while (MOVE_QUEUE.length > 0) {
        const move = MOVE_QUEUE.shift()
        renderRemoved(move)
    }
    document.querySelector('.floating-modal')?.classList.add('hidden')
}

function makeMove(row, col) {
    if(!WINNER) {
        const last = MOVE_QUEUE[MOVE_QUEUE.length - 1]
        const player = last?.player === 'X' ? 'O' : 'X'
        const move = { pos: [row, col], player }
        MOVE_QUEUE.push(move)
        const removed = MOVE_QUEUE.length > 6 ? MOVE_QUEUE.shift() : null
        renderAdded(move)
        renderRemoved(removed)
        const win = isWinningMove(MOVE_QUEUE)
        if (win) {
            renderWin(MOVE_QUEUE)
            WINNER = true
        }
    }
}

function isWinningMove(queue) {
    let result = false
    const player = queue[queue.length - 1]?.player
    const playerMoves = queue.filter(m => m.player === player)
    if (playerMoves.length >= 3) {
        const [a, b, c] = playerMoves.slice(-3).map(m => m.pos)
        const [x1, y1] = a
        const [x2, y2] = b
        const [x3, y3] = c
        const area = x1 * (y2 - y3) + x2 * (y3 - y1) + x3 * (y1 - y2)
        result = area === 0
    }
    return result
}

function renderWin(queue) {
    const player = queue[queue.length - 1]?.player
    const playerMoves = queue.filter(m => m.player === player)
    if (playerMoves.length === 3) {
        queue.forEach(m => {
            if (m.player === player) {
                const [row, col] = m.pos
                const cell = document.querySelectorAll('.tic-back-toe-cell')[row * 3 + col]
                cell.classList.add('win')
            }
        })
        document.querySelector('.floating-modal')?.classList.remove('hidden')
    }
}

function renderAdded(move) {
    if (move) {
        const [row, col] = move.pos
        const cell = document.querySelectorAll('.tic-back-toe-cell')[row * 3 + col]
        const drawClass = move.player === 'X' ? 'draw-x' : 'draw-o'
        cell.classList.add(drawClass, 'tile-appear')
        requestAnimationFrame(() => cell.classList.remove('tile-appear'))
    }
}

function renderRemoved(move) {
    if (move) {
        const [row, col] = move.pos
        const cell = document.querySelectorAll('.tic-back-toe-cell')[row * 3 + col]
        cell.classList.add('disappear')
        setTimeout(() => {
            cell.classList.remove('disappear', 'draw-x', 'draw-o', 'win')
        }, 200)
    }
}


/*
    Main
*/
document.addEventListener('DOMContentLoaded', async () => {
    const cells = document.querySelectorAll('.tic-back-toe-cell')

    cells.forEach((cell, i) => {
        cell.addEventListener('click', () => {
            const row = Math.floor(i / 3)
            const col = i % 3

            if (!MOVE_QUEUE.some(m => m.pos[0] === row && m.pos[1] === col)) {
                makeMove(row, col)
            }
        })
    })
    document.getElementById('playAgainButton')?.addEventListener('click', resetGame)

    await getUserAsync();
})

