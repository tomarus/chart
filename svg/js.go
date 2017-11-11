package svg

const js = `
let active, mkx, mkx2, mky, mky2, mks, mkt, loc, selx, sely, seltxt="", mavtxt="", mav=0, selmode=0
let dopt = {year: "numeric", month: "2-digit", day: "2-digit", hour: "2-digit", minute: "2-digit", hour12: false}
window.onload = init
document.addEventListener('load', init)
function id(n) { return 'path'+(n+1) }
function idb(n) { return 'path'+(n+1)+'_b' }
function init() {
	data.forEach((d, i) => {
		document.getElementById(idb(i)).onclick = () => { click(i); selmode = 0 }
	})
	document.getElementById('mabut').onclick = () => { maclick(); selmode = 0 }
	render(0)
	mkx = document.getElementById('markerx')
	mky = document.getElementById('markery')
	mkx2 = document.getElementById('markerx2')
	mky2 = document.getElementById('markery2')
	mkt = document.getElementById('markertext')
	mks = document.getElementById('markersel')
	handlemouse()
}
function handlemouse() {
	let svg = document.querySelector('svg')
	let pt = svg.createSVGPoint()
	svg.addEventListener('mousedown', function(evt) { 
		selmode++
		if (selmode > 2) { 
			selmode = 0
			markerpos(svg, pt, evt)
		}
		selx = loc.x
		sely = loc.y
		evt.preventDefault()
	})
	svg.addEventListener('mouseup', function(evt) { 
		if (selmode == 1) selmode = 2
		evt.preventDefault()
	})
	svg.addEventListener('mousemove', function(evt) {
		markerpos(svg, pt, evt)
	})
}
function markerpos(svg, pt,evt) {
	pt.x = evt.clientX
	pt.y = evt.clientY
	loc = pt.matrixTransform(svg.getScreenCTM().inverse())
	if (loc.x<mx || loc.x>=w+mx || loc.y<my || loc.y>=h+my) {
		marker('hidden', 0, 0)
	} else {
		marker('visible', loc.x, loc.y)
	}
	status()
	evt.preventDefault()
}
function mkvis(v) {
	mkx2.style.visibility = v
	mky2.style.visibility = v
	mks.style.visibility = v
}
function marker(s, x, y) {
	if (selmode == 2) {
		return
	}
	mkx.style.visibility = s
	mky.style.visibility = s
	mkvis(selmode?s:'hidden')
	if (s!=='visible') {
		seltxt = ""
		return
	}
	let sx, sy, dx, dy
	if (selmode) {
		mkx2.setAttribute('x1', x)
		mkx2.setAttribute('x2', x)
		mky2.setAttribute('y1', y)
		mky2.setAttribute('y2', y)
		sx = selx>x?x:selx
		sy = sely>y?y:sely
		dx = Math.abs(selx-x)
		dy = Math.abs(sely-y)
		mks.setAttribute('x', sx)
		mks.setAttribute('y', sy)
		mks.setAttribute('width', dx)
		mks.setAttribute('height', dy)
	} else {
		mkx.setAttribute('x1', x)
		mkx.setAttribute('x2', x)
		mky.setAttribute('y1', y)
		mky.setAttribute('y2', y)
		selx=w
	}
	let px = x-mx
	let py = y-my
	let v = data[active||0].fmax - data[active||0].fmax / h * py
	if (selx > 0 && selx != w) {
		d = new Date((end - start) / w * Math.min(px, selx-mx||0) + start)
	} else {
		d = new Date(((end-start) / w * px) + start)
	}
	seltxt = d.toLocaleTimeString('nl-NL', dopt)
	if (selmode) {
		let v2 = data[active||0].fmax / h * dy
		let t = (end-start) / w * dx / 1000
		seltxt += ' Len: ' + fmtime(t) + ' Delta-Y: ' + fmt(Math.abs(v2))
	} else {
		seltxt += ' Y:' + fmt(v)
	}
}
function status() {
	mkt.innerHTML = seltxt + (mavtxt !== "" ? " " + mavtxt : "")
}
function fmt(b) {
	if (b < 1000000) return b.toFixed()
	let sizes = ['', 'K', 'M', 'G', 'T', 'P']
	let i = Math.floor(Math.log(b) / Math.log(1000))
	return parseFloat((b / Math.pow(1000, i))).toFixed(3) + '' + sizes[i]
}
function fmtime(t) {
	let d = Math.floor(t/86400)
	let h = Math.floor(t/3600)%24
	let m = Math.floor(t/60)%60
	return (d>0?d+'d ':'')+(h>0?h+'h ':'')+(m>0?m+'m':'')
}
function click(n) {
	if (active === n) {
		styles('visible', 1)
		render(0)
		active = undefined
		return
	}
	styles('hidden', 0.20)
	style(n, 'visible', 1)
	render(n)
	active = n
}
function style(n, v, o) {
	document.getElementById(id(n)).style.visibility = v
	document.getElementById(idb(n)).style.fillOpacity = o
}
function styles(v, o) {
	data.forEach((d, i) => {
		style(i, v, o)
	})
}
function scale(n) {
	let c = document.getElementById('ygrid').children
	for (i=0; i<c.length; i++) {
		c[i].children[0].innerHTML = data[n].scale[i]
	}
}
function norm(v, max, fmax) {
	let a = fmax / max
	let b = fmax - a*max
	return v * a + b
}
function render(n) {
	scale(n)
	data.forEach((d, i) => {
		window[data[i].type](i, h, i === n ? h : data[i].max)
	})
	ma(n)
}
function area(n, max, fmax) {
	graph('area', '', n, max, fmax)
}
function line(n, max, fmax) {
	let v0 = norm(data[n].values[0], max, fmax)
	graph('line', 'M0,'+(h-v0), n, max ,fmax)
}
function graph(t, p, n, max, fmax) {
	for (let i=0; i<Math.min(w, data[n].values.length); i++) {
		let v = norm(data[n].values[i], max, fmax)
		// if (t === 'line') {
		// 	p += 'L'+i+','+(h-v)
		// } else {
			p += 'M'+i+','+(h-v)+'v'+v
		// }
	}
	document.getElementById(id(n)).firstElementChild.setAttribute('d', p)
}
function valof(n, pos, off, max) {
	if (pos+off<0 || pos+off>=max) {
		return -1
	}
	return data[n].values[pos+off]
}
function findma(n, pos, size) {
	const minValue = 1
	if (size == w) {
		let v = 0, tw = 0
		for (let j=0; j<size; j++) {
			if (data[n].values[j]>=minValue) {
				v += data[n].values[j]
				tw++
			}
		}
		return v/tw
	}
	let v = data[n].values[pos]
	let dx = data[n].values.length
	let totw = 0
	for (let j=0; j<size; j++) {
		let wx = size-j-1
		let q = valof(n, pos, j+1, dx)
		if (q>=minValue) {
			totw += wx
			v += q*wx
		}
		q = valof(n, pos, -j-1, dx)
		if (q>=minValue) {
			totw += wx
			v += q*wx
		}
	}
	v /= totw+1
	return v
}
function ma(n) {
	if (mav===0) {
		document.getElementById('ma').firstElementChild.setAttribute('d', 'M0,0')
		return
	}
	let smooth = 1<<mav
	if (smooth>w) {
		smooth = w
	}
	let v0 = norm(findma(n, 0, smooth), h, h)
	let p = 'M0,'+(h-v0)
	for (let i=0; i<Math.min(w, data[n].values.length); i++) {
		let v = norm(findma(n, i, smooth), h, h)
		p += 'L'+i+','+(h-v)
	}
	document.getElementById('ma').firstElementChild.setAttribute('d', p)
}
function maclick() {
	if (mav === w || 1<<mav > w) {
		mav = -1
	}
	mav++
	ma(active||0)
	mavtxt = ""
	if (mav>0) {
		mavtxt = "WMA:" + (1<<mav > w ? "all" : 1<<mav)
	}
	status()
}
`
