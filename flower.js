const c = document.getElementById("myCanvas");
const ctx = c.getContext("2d");
let start;
let grid = {
    dat: [],
    w: 0,
    h: 0,
    ox: 0,
    oy: 0,
}

ctx.canvas.width = window.innerWidth;
ctx.canvas.height = window.innerHeight;
window.addEventListener('resize', function () {
    ctx.canvas.width = window.innerWidth;
    ctx.canvas.height = window.innerHeight;
    console.log(document.width, document.height);
}, false);

function url(s) {
    var l = window.location;
    return ((l.protocol === "https:") ? "wss://" : "ws://") + l.hostname + (((l.port != 80) && (l.port != 443)) ? ":" + l.port : "") + l.pathname + s;
}

const socket = new WebSocket(url('ws'));
socket.onopen = ev => {
    socket.send(JSON.stringify({
        Color: "yuh",
        I: 0,
        J: 4
    }))
}
socket.onmessage = ev => {
    console.log(ev.data);
    const msg = JSON.parse(ev.data);
    msg.forEach(m => {
        ctx.fillStyle = m.Color
        ctx.fillRect(m.I - 5, m.J - 5, 10, 10);
    });

}

function getGrid({dat, w, h, ox, oy}, i, j) {
    const ii = i + ox;
    const jj = j + oy;
    if (ii < 0 || jj < 0 || ii >= w || jj >= h) {
        return 0;
    }
    const idx = jj * width + ii;
    return dat[idx];
}

let x = 0;
let y = 0;
let isDrawing = false;

c.addEventListener('mousedown', tStart);
c.addEventListener('touchstart', tStart);
c.addEventListener('touchmove', tMove);
c.addEventListener('touchend', tEnd);
c.addEventListener("mousemove",tMove);
c.addEventListener("mouseup", tEnd);

function tStart(e) {
    isDrawing = true;
    console.log("down");
}

function tMove(e) {
        rainbow++;
        if (rainbow >= 360) {
            rainbow = 0;
        }
        if (isDrawing) {
            // ctx.fillStyle = "red";
            //ctx.fillStyle = 'hsl('+ 360*Math.random() +',100%,50%)';
            const msg = {
                Color: 'hsl(' + rainbow + ',100%,50%)',
                I: e.offsetX,
                J: e.offsetY,
            };
            socket.send(JSON.stringify(msg))

            //ctx.fillRect(e.offsetX - 5, e.offsetY -5, 10, 10);
        }
        console.log("move", e.offsetX, e.offsetY, isDrawing);
}

function tEnd(e) {
    isDrawing = false;
}


let rainbow = 0;

function setInplace({dat, w, h, ox, oy}, i, j) {

}

function step(ts) {
    if (!start) {
        start = ts;
    }

    ctx.fillStyle = "purple";
    ctx.fillRect(0, 0, c.width / 2, c.height);
    ctx.fillStyle = "#552200";
    ctx.fillRect(c.width / 2, 0, c.width / 2, c.height);
    //ctx.fill();
    //window.requestAnimationFrame(step);
}

step();
//window.requestAnimationFrame(step);
