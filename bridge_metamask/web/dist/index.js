(() => {
  var __create = Object.create;
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getProtoOf = Object.getPrototypeOf;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __commonJS = (cb, mod2) => function __require() {
    return mod2 || (0, cb[__getOwnPropNames(cb)[0]])((mod2 = { exports: {} }).exports, mod2), mod2.exports;
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toESM = (mod2, isNodeMode, target) => (target = mod2 != null ? __create(__getProtoOf(mod2)) : {}, __copyProps(
    // If the importer is in node compatibility mode or this is not an ESM
    // file that has been converted to a CommonJS file using a Babel-
    // compatible transform (i.e. "__esModule" has not been set), then set
    // "default" to the CommonJS "module.exports" for node compatibility.
    isNodeMode || !mod2 || !mod2.__esModule ? __defProp(target, "default", { value: mod2, enumerable: true }) : target,
    mod2
  ));

  // (disabled):crypto
  var require_crypto = __commonJS({
    "(disabled):crypto"() {
    }
  });

  // node_modules/@metamask/detect-provider/dist/index.js
  var require_dist = __commonJS({
    "node_modules/@metamask/detect-provider/dist/index.js"(exports, module) {
      "use strict";
      function detectEthereumProvider({ mustBeMetaMask = false, silent = false, timeout = 3e3 } = {}) {
        _validateInputs();
        let handled = false;
        return new Promise((resolve) => {
          if (window.ethereum) {
            handleEthereum();
          } else {
            window.addEventListener("ethereum#initialized", handleEthereum, { once: true });
            setTimeout(() => {
              handleEthereum();
            }, timeout);
          }
          function handleEthereum() {
            if (handled) {
              return;
            }
            handled = true;
            window.removeEventListener("ethereum#initialized", handleEthereum);
            const { ethereum } = window;
            if (ethereum && (!mustBeMetaMask || ethereum.isMetaMask)) {
              resolve(ethereum);
            } else {
              const message = mustBeMetaMask && ethereum ? "Non-MetaMask window.ethereum detected." : "Unable to detect window.ethereum.";
              !silent && console.error("@metamask/detect-provider:", message);
              resolve(null);
            }
          }
        });
        function _validateInputs() {
          if (typeof mustBeMetaMask !== "boolean") {
            throw new Error(`@metamask/detect-provider: Expected option 'mustBeMetaMask' to be a boolean.`);
          }
          if (typeof silent !== "boolean") {
            throw new Error(`@metamask/detect-provider: Expected option 'silent' to be a boolean.`);
          }
          if (typeof timeout !== "number") {
            throw new Error(`@metamask/detect-provider: Expected option 'timeout' to be a number.`);
          }
        }
      }
      module.exports = detectEthereumProvider;
    }
  });

  // node_modules/preact/dist/preact.module.js
  var n;
  var l;
  var u;
  var t;
  var i;
  var o;
  var r;
  var f;
  var e;
  var c = {};
  var s = [];
  var a = /acit|ex(?:s|g|n|p|$)|rph|grid|ows|mnc|ntw|ine[ch]|zoo|^ord|itera/i;
  var h = Array.isArray;
  function v(n4, l4) {
    for (var u4 in l4)
      n4[u4] = l4[u4];
    return n4;
  }
  function p(n4) {
    var l4 = n4.parentNode;
    l4 && l4.removeChild(n4);
  }
  function y(l4, u4, t4) {
    var i4, o4, r4, f4 = {};
    for (r4 in u4)
      "key" == r4 ? i4 = u4[r4] : "ref" == r4 ? o4 = u4[r4] : f4[r4] = u4[r4];
    if (arguments.length > 2 && (f4.children = arguments.length > 3 ? n.call(arguments, 2) : t4), "function" == typeof l4 && null != l4.defaultProps)
      for (r4 in l4.defaultProps)
        void 0 === f4[r4] && (f4[r4] = l4.defaultProps[r4]);
    return d(l4, f4, i4, o4, null);
  }
  function d(n4, t4, i4, o4, r4) {
    var f4 = { type: n4, props: t4, key: i4, ref: o4, __k: null, __: null, __b: 0, __e: null, __d: void 0, __c: null, __h: null, constructor: void 0, __v: null == r4 ? ++u : r4 };
    return null == r4 && null != l.vnode && l.vnode(f4), f4;
  }
  function k(n4) {
    return n4.children;
  }
  function b(n4, l4) {
    this.props = n4, this.context = l4;
  }
  function g(n4, l4) {
    if (null == l4)
      return n4.__ ? g(n4.__, n4.__.__k.indexOf(n4) + 1) : null;
    for (var u4; l4 < n4.__k.length; l4++)
      if (null != (u4 = n4.__k[l4]) && null != u4.__e)
        return u4.__e;
    return "function" == typeof n4.type ? g(n4) : null;
  }
  function m(n4) {
    var l4, u4;
    if (null != (n4 = n4.__) && null != n4.__c) {
      for (n4.__e = n4.__c.base = null, l4 = 0; l4 < n4.__k.length; l4++)
        if (null != (u4 = n4.__k[l4]) && null != u4.__e) {
          n4.__e = n4.__c.base = u4.__e;
          break;
        }
      return m(n4);
    }
  }
  function w(n4) {
    (!n4.__d && (n4.__d = true) && i.push(n4) && !x.__r++ || o !== l.debounceRendering) && ((o = l.debounceRendering) || r)(x);
  }
  function x() {
    var n4, l4, u4, t4, o4, r4, e4, c4, s4;
    for (i.sort(f); n4 = i.shift(); )
      n4.__d && (l4 = i.length, t4 = void 0, o4 = void 0, r4 = void 0, c4 = (e4 = (u4 = n4).__v).__e, (s4 = u4.__P) && (t4 = [], o4 = [], (r4 = v({}, e4)).__v = e4.__v + 1, L(s4, e4, r4, u4.__n, void 0 !== s4.ownerSVGElement, null != e4.__h ? [c4] : null, t4, null == c4 ? g(e4) : c4, e4.__h, o4), M(t4, e4, o4), e4.__e != c4 && m(e4)), i.length > l4 && i.sort(f));
    x.__r = 0;
  }
  function P(n4, l4, u4, t4, i4, o4, r4, f4, e4, a4, v3) {
    var p4, y2, _, b5, m4, w3, x2, P2, C, H2 = 0, I2 = t4 && t4.__k || s, T3 = I2.length, j4 = T3, z3 = l4.length;
    for (u4.__k = [], p4 = 0; p4 < z3; p4++)
      null != (b5 = u4.__k[p4] = null == (b5 = l4[p4]) || "boolean" == typeof b5 || "function" == typeof b5 ? null : "string" == typeof b5 || "number" == typeof b5 || "bigint" == typeof b5 ? d(null, b5, null, null, b5) : h(b5) ? d(k, { children: b5 }, null, null, null) : b5.__b > 0 ? d(b5.type, b5.props, b5.key, b5.ref ? b5.ref : null, b5.__v) : b5) ? (b5.__ = u4, b5.__b = u4.__b + 1, -1 === (P2 = A(b5, I2, x2 = p4 + H2, j4)) ? _ = c : (_ = I2[P2] || c, I2[P2] = void 0, j4--), L(n4, b5, _, i4, o4, r4, f4, e4, a4, v3), m4 = b5.__e, (y2 = b5.ref) && _.ref != y2 && (_.ref && O(_.ref, null, b5), v3.push(y2, b5.__c || m4, b5)), null != m4 && (null == w3 && (w3 = m4), (C = _ === c || null === _.__v) ? -1 == P2 && H2-- : P2 !== x2 && (P2 === x2 + 1 ? H2++ : P2 > x2 ? j4 > z3 - x2 ? H2 += P2 - x2 : H2-- : H2 = P2 < x2 && P2 == x2 - 1 ? P2 - x2 : 0), x2 = p4 + H2, "function" != typeof b5.type || P2 === x2 && _.__k !== b5.__k ? "function" == typeof b5.type || P2 === x2 && !C ? void 0 !== b5.__d ? (e4 = b5.__d, b5.__d = void 0) : e4 = m4.nextSibling : e4 = S(n4, m4, e4) : e4 = $(b5, e4, n4), "function" == typeof u4.type && (u4.__d = e4))) : (_ = I2[p4]) && null == _.key && _.__e && (_.__e == e4 && (e4 = g(_)), q(_, _, false), I2[p4] = null);
    for (u4.__e = w3, p4 = T3; p4--; )
      null != I2[p4] && ("function" == typeof u4.type && null != I2[p4].__e && I2[p4].__e == u4.__d && (u4.__d = I2[p4].__e.nextSibling), q(I2[p4], I2[p4]));
  }
  function $(n4, l4, u4) {
    for (var t4, i4 = n4.__k, o4 = 0; i4 && o4 < i4.length; o4++)
      (t4 = i4[o4]) && (t4.__ = n4, l4 = "function" == typeof t4.type ? $(t4, l4, u4) : S(u4, t4.__e, l4));
    return l4;
  }
  function S(n4, l4, u4) {
    return null == u4 || u4.parentNode !== n4 ? n4.insertBefore(l4, null) : l4 == u4 && null != l4.parentNode || n4.insertBefore(l4, u4), l4.nextSibling;
  }
  function A(n4, l4, u4, t4) {
    var i4 = n4.key, o4 = n4.type, r4 = u4 - 1, f4 = u4 + 1, e4 = l4[u4];
    if (null === e4 || e4 && i4 == e4.key && o4 === e4.type)
      return u4;
    if (t4 > (null != e4 ? 1 : 0))
      for (; r4 >= 0 || f4 < l4.length; ) {
        if (r4 >= 0) {
          if ((e4 = l4[r4]) && i4 == e4.key && o4 === e4.type)
            return r4;
          r4--;
        }
        if (f4 < l4.length) {
          if ((e4 = l4[f4]) && i4 == e4.key && o4 === e4.type)
            return f4;
          f4++;
        }
      }
    return -1;
  }
  function H(n4, l4, u4, t4, i4) {
    var o4;
    for (o4 in u4)
      "children" === o4 || "key" === o4 || o4 in l4 || T(n4, o4, null, u4[o4], t4);
    for (o4 in l4)
      i4 && "function" != typeof l4[o4] || "children" === o4 || "key" === o4 || "value" === o4 || "checked" === o4 || u4[o4] === l4[o4] || T(n4, o4, l4[o4], u4[o4], t4);
  }
  function I(n4, l4, u4) {
    "-" === l4[0] ? n4.setProperty(l4, null == u4 ? "" : u4) : n4[l4] = null == u4 ? "" : "number" != typeof u4 || a.test(l4) ? u4 : u4 + "px";
  }
  function T(n4, l4, u4, t4, i4) {
    var o4;
    n:
      if ("style" === l4)
        if ("string" == typeof u4)
          n4.style.cssText = u4;
        else {
          if ("string" == typeof t4 && (n4.style.cssText = t4 = ""), t4)
            for (l4 in t4)
              u4 && l4 in u4 || I(n4.style, l4, "");
          if (u4)
            for (l4 in u4)
              t4 && u4[l4] === t4[l4] || I(n4.style, l4, u4[l4]);
        }
      else if ("o" === l4[0] && "n" === l4[1])
        o4 = l4 !== (l4 = l4.replace(/(PointerCapture)$|Capture$/, "$1")), l4 = l4.toLowerCase() in n4 ? l4.toLowerCase().slice(2) : l4.slice(2), n4.l || (n4.l = {}), n4.l[l4 + o4] = u4, u4 ? t4 || n4.addEventListener(l4, o4 ? z : j, o4) : n4.removeEventListener(l4, o4 ? z : j, o4);
      else if ("dangerouslySetInnerHTML" !== l4) {
        if (i4)
          l4 = l4.replace(/xlink(H|:h)/, "h").replace(/sName$/, "s");
        else if ("width" !== l4 && "height" !== l4 && "href" !== l4 && "list" !== l4 && "form" !== l4 && "tabIndex" !== l4 && "download" !== l4 && "rowSpan" !== l4 && "colSpan" !== l4 && l4 in n4)
          try {
            n4[l4] = null == u4 ? "" : u4;
            break n;
          } catch (n5) {
          }
        "function" == typeof u4 || (null == u4 || false === u4 && "-" !== l4[4] ? n4.removeAttribute(l4) : n4.setAttribute(l4, u4));
      }
  }
  function j(n4) {
    return this.l[n4.type + false](l.event ? l.event(n4) : n4);
  }
  function z(n4) {
    return this.l[n4.type + true](l.event ? l.event(n4) : n4);
  }
  function L(n4, u4, t4, i4, o4, r4, f4, e4, c4, s4) {
    var a4, p4, y2, d4, _, g4, m4, w3, x2, $2, C, S2, A2, H2, I2, T3 = u4.type;
    if (void 0 !== u4.constructor)
      return null;
    null != t4.__h && (c4 = t4.__h, e4 = u4.__e = t4.__e, u4.__h = null, r4 = [e4]), (a4 = l.__b) && a4(u4);
    n:
      if ("function" == typeof T3)
        try {
          if (w3 = u4.props, x2 = (a4 = T3.contextType) && i4[a4.__c], $2 = a4 ? x2 ? x2.props.value : a4.__ : i4, t4.__c ? m4 = (p4 = u4.__c = t4.__c).__ = p4.__E : ("prototype" in T3 && T3.prototype.render ? u4.__c = p4 = new T3(w3, $2) : (u4.__c = p4 = new b(w3, $2), p4.constructor = T3, p4.render = B), x2 && x2.sub(p4), p4.props = w3, p4.state || (p4.state = {}), p4.context = $2, p4.__n = i4, y2 = p4.__d = true, p4.__h = [], p4._sb = []), null == p4.__s && (p4.__s = p4.state), null != T3.getDerivedStateFromProps && (p4.__s == p4.state && (p4.__s = v({}, p4.__s)), v(p4.__s, T3.getDerivedStateFromProps(w3, p4.__s))), d4 = p4.props, _ = p4.state, p4.__v = u4, y2)
            null == T3.getDerivedStateFromProps && null != p4.componentWillMount && p4.componentWillMount(), null != p4.componentDidMount && p4.__h.push(p4.componentDidMount);
          else {
            if (null == T3.getDerivedStateFromProps && w3 !== d4 && null != p4.componentWillReceiveProps && p4.componentWillReceiveProps(w3, $2), !p4.__e && (null != p4.shouldComponentUpdate && false === p4.shouldComponentUpdate(w3, p4.__s, $2) || u4.__v === t4.__v)) {
              for (u4.__v !== t4.__v && (p4.props = w3, p4.state = p4.__s, p4.__d = false), u4.__e = t4.__e, u4.__k = t4.__k, u4.__k.forEach(function(n5) {
                n5 && (n5.__ = u4);
              }), C = 0; C < p4._sb.length; C++)
                p4.__h.push(p4._sb[C]);
              p4._sb = [], p4.__h.length && f4.push(p4);
              break n;
            }
            null != p4.componentWillUpdate && p4.componentWillUpdate(w3, p4.__s, $2), null != p4.componentDidUpdate && p4.__h.push(function() {
              p4.componentDidUpdate(d4, _, g4);
            });
          }
          if (p4.context = $2, p4.props = w3, p4.__P = n4, p4.__e = false, S2 = l.__r, A2 = 0, "prototype" in T3 && T3.prototype.render) {
            for (p4.state = p4.__s, p4.__d = false, S2 && S2(u4), a4 = p4.render(p4.props, p4.state, p4.context), H2 = 0; H2 < p4._sb.length; H2++)
              p4.__h.push(p4._sb[H2]);
            p4._sb = [];
          } else
            do {
              p4.__d = false, S2 && S2(u4), a4 = p4.render(p4.props, p4.state, p4.context), p4.state = p4.__s;
            } while (p4.__d && ++A2 < 25);
          p4.state = p4.__s, null != p4.getChildContext && (i4 = v(v({}, i4), p4.getChildContext())), y2 || null == p4.getSnapshotBeforeUpdate || (g4 = p4.getSnapshotBeforeUpdate(d4, _)), P(n4, h(I2 = null != a4 && a4.type === k && null == a4.key ? a4.props.children : a4) ? I2 : [I2], u4, t4, i4, o4, r4, f4, e4, c4, s4), p4.base = u4.__e, u4.__h = null, p4.__h.length && f4.push(p4), m4 && (p4.__E = p4.__ = null);
        } catch (n5) {
          u4.__v = null, (c4 || null != r4) && (u4.__e = e4, u4.__h = !!c4, r4[r4.indexOf(e4)] = null), l.__e(n5, u4, t4);
        }
      else
        null == r4 && u4.__v === t4.__v ? (u4.__k = t4.__k, u4.__e = t4.__e) : u4.__e = N(t4.__e, u4, t4, i4, o4, r4, f4, c4, s4);
    (a4 = l.diffed) && a4(u4);
  }
  function M(n4, u4, t4) {
    for (var i4 = 0; i4 < t4.length; i4++)
      O(t4[i4], t4[++i4], t4[++i4]);
    l.__c && l.__c(u4, n4), n4.some(function(u5) {
      try {
        n4 = u5.__h, u5.__h = [], n4.some(function(n5) {
          n5.call(u5);
        });
      } catch (n5) {
        l.__e(n5, u5.__v);
      }
    });
  }
  function N(l4, u4, t4, i4, o4, r4, f4, e4, s4) {
    var a4, v3, y2, d4 = t4.props, _ = u4.props, k3 = u4.type, b5 = 0;
    if ("svg" === k3 && (o4 = true), null != r4) {
      for (; b5 < r4.length; b5++)
        if ((a4 = r4[b5]) && "setAttribute" in a4 == !!k3 && (k3 ? a4.localName === k3 : 3 === a4.nodeType)) {
          l4 = a4, r4[b5] = null;
          break;
        }
    }
    if (null == l4) {
      if (null === k3)
        return document.createTextNode(_);
      l4 = o4 ? document.createElementNS("http://www.w3.org/2000/svg", k3) : document.createElement(k3, _.is && _), r4 = null, e4 = false;
    }
    if (null === k3)
      d4 === _ || e4 && l4.data === _ || (l4.data = _);
    else {
      if (r4 = r4 && n.call(l4.childNodes), v3 = (d4 = t4.props || c).dangerouslySetInnerHTML, y2 = _.dangerouslySetInnerHTML, !e4) {
        if (null != r4)
          for (d4 = {}, b5 = 0; b5 < l4.attributes.length; b5++)
            d4[l4.attributes[b5].name] = l4.attributes[b5].value;
        (y2 || v3) && (y2 && (v3 && y2.__html == v3.__html || y2.__html === l4.innerHTML) || (l4.innerHTML = y2 && y2.__html || ""));
      }
      if (H(l4, _, d4, o4, e4), y2)
        u4.__k = [];
      else if (P(l4, h(b5 = u4.props.children) ? b5 : [b5], u4, t4, i4, o4 && "foreignObject" !== k3, r4, f4, r4 ? r4[0] : t4.__k && g(t4, 0), e4, s4), null != r4)
        for (b5 = r4.length; b5--; )
          null != r4[b5] && p(r4[b5]);
      e4 || ("value" in _ && void 0 !== (b5 = _.value) && (b5 !== l4.value || "progress" === k3 && !b5 || "option" === k3 && b5 !== d4.value) && T(l4, "value", b5, d4.value, false), "checked" in _ && void 0 !== (b5 = _.checked) && b5 !== l4.checked && T(l4, "checked", b5, d4.checked, false));
    }
    return l4;
  }
  function O(n4, u4, t4) {
    try {
      "function" == typeof n4 ? n4(u4) : n4.current = u4;
    } catch (n5) {
      l.__e(n5, t4);
    }
  }
  function q(n4, u4, t4) {
    var i4, o4;
    if (l.unmount && l.unmount(n4), (i4 = n4.ref) && (i4.current && i4.current !== n4.__e || O(i4, null, u4)), null != (i4 = n4.__c)) {
      if (i4.componentWillUnmount)
        try {
          i4.componentWillUnmount();
        } catch (n5) {
          l.__e(n5, u4);
        }
      i4.base = i4.__P = null, n4.__c = void 0;
    }
    if (i4 = n4.__k)
      for (o4 = 0; o4 < i4.length; o4++)
        i4[o4] && q(i4[o4], u4, t4 || "function" != typeof n4.type);
    t4 || null == n4.__e || p(n4.__e), n4.__ = n4.__e = n4.__d = void 0;
  }
  function B(n4, l4, u4) {
    return this.constructor(n4, u4);
  }
  function D(u4, t4, i4) {
    var o4, r4, f4, e4;
    l.__ && l.__(u4, t4), r4 = (o4 = "function" == typeof i4) ? null : i4 && i4.__k || t4.__k, f4 = [], e4 = [], L(t4, u4 = (!o4 && i4 || t4).__k = y(k, null, [u4]), r4 || c, c, void 0 !== t4.ownerSVGElement, !o4 && i4 ? [i4] : r4 ? null : t4.firstChild ? n.call(t4.childNodes) : null, f4, !o4 && i4 ? i4 : r4 ? r4.__e : t4.firstChild, o4, e4), M(f4, u4, e4);
  }
  n = s.slice, l = { __e: function(n4, l4, u4, t4) {
    for (var i4, o4, r4; l4 = l4.__; )
      if ((i4 = l4.__c) && !i4.__)
        try {
          if ((o4 = i4.constructor) && null != o4.getDerivedStateFromError && (i4.setState(o4.getDerivedStateFromError(n4)), r4 = i4.__d), null != i4.componentDidCatch && (i4.componentDidCatch(n4, t4 || {}), r4 = i4.__d), r4)
            return i4.__E = i4;
        } catch (l5) {
          n4 = l5;
        }
    throw n4;
  } }, u = 0, t = function(n4) {
    return null != n4 && void 0 === n4.constructor;
  }, b.prototype.setState = function(n4, l4) {
    var u4;
    u4 = null != this.__s && this.__s !== this.state ? this.__s : this.__s = v({}, this.state), "function" == typeof n4 && (n4 = n4(v({}, u4), this.props)), n4 && v(u4, n4), null != n4 && this.__v && (l4 && this._sb.push(l4), w(this));
  }, b.prototype.forceUpdate = function(n4) {
    this.__v && (this.__e = true, n4 && this.__h.push(n4), w(this));
  }, b.prototype.render = k, i = [], r = "function" == typeof Promise ? Promise.prototype.then.bind(Promise.resolve()) : setTimeout, f = function(n4, l4) {
    return n4.__v.__b - l4.__v.__b;
  }, x.__r = 0, e = 0;

  // node_modules/preact/hooks/dist/hooks.module.js
  var t2;
  var r2;
  var u2;
  var i2;
  var o2 = 0;
  var f2 = [];
  var c2 = [];
  var e2 = l.__b;
  var a2 = l.__r;
  var v2 = l.diffed;
  var l2 = l.__c;
  var m2 = l.unmount;
  function d2(t4, u4) {
    l.__h && l.__h(r2, t4, o2 || u4), o2 = 0;
    var i4 = r2.__H || (r2.__H = { __: [], __h: [] });
    return t4 >= i4.__.length && i4.__.push({ __V: c2 }), i4.__[t4];
  }
  function h2(n4) {
    return o2 = 1, s2(B2, n4);
  }
  function s2(n4, u4, i4) {
    var o4 = d2(t2++, 2);
    if (o4.t = n4, !o4.__c && (o4.__ = [i4 ? i4(u4) : B2(void 0, u4), function(n5) {
      var t4 = o4.__N ? o4.__N[0] : o4.__[0], r4 = o4.t(t4, n5);
      t4 !== r4 && (o4.__N = [r4, o4.__[1]], o4.__c.setState({}));
    }], o4.__c = r2, !r2.u)) {
      var f4 = function(n5, t4, r4) {
        if (!o4.__c.__H)
          return true;
        var u5 = o4.__c.__H.__.filter(function(n6) {
          return n6.__c;
        });
        if (u5.every(function(n6) {
          return !n6.__N;
        }))
          return !c4 || c4.call(this, n5, t4, r4);
        var i5 = false;
        return u5.forEach(function(n6) {
          if (n6.__N) {
            var t5 = n6.__[0];
            n6.__ = n6.__N, n6.__N = void 0, t5 !== n6.__[0] && (i5 = true);
          }
        }), !(!i5 && o4.__c.props === n5) && (!c4 || c4.call(this, n5, t4, r4));
      };
      r2.u = true;
      var c4 = r2.shouldComponentUpdate, e4 = r2.componentWillUpdate;
      r2.componentWillUpdate = function(n5, t4, r4) {
        if (this.__e) {
          var u5 = c4;
          c4 = void 0, f4(n5, t4, r4), c4 = u5;
        }
        e4 && e4.call(this, n5, t4, r4);
      }, r2.shouldComponentUpdate = f4;
    }
    return o4.__N || o4.__;
  }
  function p2(u4, i4) {
    var o4 = d2(t2++, 3);
    !l.__s && z2(o4.__H, i4) && (o4.__ = u4, o4.i = i4, r2.__H.__h.push(o4));
  }
  function F(n4, r4) {
    var u4 = d2(t2++, 7);
    return z2(u4.__H, r4) ? (u4.__V = n4(), u4.i = r4, u4.__h = n4, u4.__V) : u4.__;
  }
  function T2(n4, t4) {
    return o2 = 8, F(function() {
      return n4;
    }, t4);
  }
  function b2() {
    for (var t4; t4 = f2.shift(); )
      if (t4.__P && t4.__H)
        try {
          t4.__H.__h.forEach(k2), t4.__H.__h.forEach(w2), t4.__H.__h = [];
        } catch (r4) {
          t4.__H.__h = [], l.__e(r4, t4.__v);
        }
  }
  l.__b = function(n4) {
    r2 = null, e2 && e2(n4);
  }, l.__r = function(n4) {
    a2 && a2(n4), t2 = 0;
    var i4 = (r2 = n4.__c).__H;
    i4 && (u2 === r2 ? (i4.__h = [], r2.__h = [], i4.__.forEach(function(n5) {
      n5.__N && (n5.__ = n5.__N), n5.__V = c2, n5.__N = n5.i = void 0;
    })) : (i4.__h.forEach(k2), i4.__h.forEach(w2), i4.__h = [], t2 = 0)), u2 = r2;
  }, l.diffed = function(t4) {
    v2 && v2(t4);
    var o4 = t4.__c;
    o4 && o4.__H && (o4.__H.__h.length && (1 !== f2.push(o4) && i2 === l.requestAnimationFrame || ((i2 = l.requestAnimationFrame) || j2)(b2)), o4.__H.__.forEach(function(n4) {
      n4.i && (n4.__H = n4.i), n4.__V !== c2 && (n4.__ = n4.__V), n4.i = void 0, n4.__V = c2;
    })), u2 = r2 = null;
  }, l.__c = function(t4, r4) {
    r4.some(function(t5) {
      try {
        t5.__h.forEach(k2), t5.__h = t5.__h.filter(function(n4) {
          return !n4.__ || w2(n4);
        });
      } catch (u4) {
        r4.some(function(n4) {
          n4.__h && (n4.__h = []);
        }), r4 = [], l.__e(u4, t5.__v);
      }
    }), l2 && l2(t4, r4);
  }, l.unmount = function(t4) {
    m2 && m2(t4);
    var r4, u4 = t4.__c;
    u4 && u4.__H && (u4.__H.__.forEach(function(n4) {
      try {
        k2(n4);
      } catch (n5) {
        r4 = n5;
      }
    }), u4.__H = void 0, r4 && l.__e(r4, u4.__v));
  };
  var g2 = "function" == typeof requestAnimationFrame;
  function j2(n4) {
    var t4, r4 = function() {
      clearTimeout(u4), g2 && cancelAnimationFrame(t4), setTimeout(n4);
    }, u4 = setTimeout(r4, 100);
    g2 && (t4 = requestAnimationFrame(r4));
  }
  function k2(n4) {
    var t4 = r2, u4 = n4.__c;
    "function" == typeof u4 && (n4.__c = void 0, u4()), r2 = t4;
  }
  function w2(n4) {
    var t4 = r2;
    n4.__c = n4.__(), r2 = t4;
  }
  function z2(n4, t4) {
    return !n4 || n4.length !== t4.length || t4.some(function(t5, r4) {
      return t5 !== n4[r4];
    });
  }
  function B2(n4, t4) {
    return "function" == typeof t4 ? t4(n4) : t4;
  }

  // node_modules/ethers/lib.esm/_version.js
  var version = "6.7.1";

  // node_modules/ethers/lib.esm/utils/properties.js
  function checkType(value, type, name) {
    const types = type.split("|").map((t4) => t4.trim());
    for (let i4 = 0; i4 < types.length; i4++) {
      switch (type) {
        case "any":
          return;
        case "bigint":
        case "boolean":
        case "number":
        case "string":
          if (typeof value === type) {
            return;
          }
      }
    }
    const error = new Error(`invalid value for type ${type}`);
    error.code = "INVALID_ARGUMENT";
    error.argument = `value.${name}`;
    error.value = value;
    throw error;
  }
  async function resolveProperties(value) {
    const keys = Object.keys(value);
    const results = await Promise.all(keys.map((k3) => Promise.resolve(value[k3])));
    return results.reduce((accum, v3, index) => {
      accum[keys[index]] = v3;
      return accum;
    }, {});
  }
  function defineProperties(target, values, types) {
    for (let key in values) {
      let value = values[key];
      const type = types ? types[key] : null;
      if (type) {
        checkType(value, type, key);
      }
      Object.defineProperty(target, key, { enumerable: true, value, writable: false });
    }
  }

  // node_modules/ethers/lib.esm/utils/errors.js
  function stringify(value) {
    if (value == null) {
      return "null";
    }
    if (Array.isArray(value)) {
      return "[ " + value.map(stringify).join(", ") + " ]";
    }
    if (value instanceof Uint8Array) {
      const HEX = "0123456789abcdef";
      let result = "0x";
      for (let i4 = 0; i4 < value.length; i4++) {
        result += HEX[value[i4] >> 4];
        result += HEX[value[i4] & 15];
      }
      return result;
    }
    if (typeof value === "object" && typeof value.toJSON === "function") {
      return stringify(value.toJSON());
    }
    switch (typeof value) {
      case "boolean":
      case "symbol":
        return value.toString();
      case "bigint":
        return BigInt(value).toString();
      case "number":
        return value.toString();
      case "string":
        return JSON.stringify(value);
      case "object": {
        const keys = Object.keys(value);
        keys.sort();
        return "{ " + keys.map((k3) => `${stringify(k3)}: ${stringify(value[k3])}`).join(", ") + " }";
      }
    }
    return `[ COULD NOT SERIALIZE ]`;
  }
  function isError(error, code) {
    return error && error.code === code;
  }
  function isCallException(error) {
    return isError(error, "CALL_EXCEPTION");
  }
  function makeError(message, code, info) {
    {
      const details = [];
      if (info) {
        if ("message" in info || "code" in info || "name" in info) {
          throw new Error(`value will overwrite populated values: ${stringify(info)}`);
        }
        for (const key in info) {
          const value = info[key];
          details.push(key + "=" + stringify(value));
        }
      }
      details.push(`code=${code}`);
      details.push(`version=${version}`);
      if (details.length) {
        message += " (" + details.join(", ") + ")";
      }
    }
    let error;
    switch (code) {
      case "INVALID_ARGUMENT":
        error = new TypeError(message);
        break;
      case "NUMERIC_FAULT":
      case "BUFFER_OVERRUN":
        error = new RangeError(message);
        break;
      default:
        error = new Error(message);
    }
    defineProperties(error, { code });
    if (info) {
      Object.assign(error, info);
    }
    return error;
  }
  function assert(check, message, code, info) {
    if (!check) {
      throw makeError(message, code, info);
    }
  }
  function assertArgument(check, message, name, value) {
    assert(check, message, "INVALID_ARGUMENT", { argument: name, value });
  }
  function assertArgumentCount(count, expectedCount, message) {
    if (message == null) {
      message = "";
    }
    if (message) {
      message = ": " + message;
    }
    assert(count >= expectedCount, "missing arguemnt" + message, "MISSING_ARGUMENT", {
      count,
      expectedCount
    });
    assert(count <= expectedCount, "too many arguemnts" + message, "UNEXPECTED_ARGUMENT", {
      count,
      expectedCount
    });
  }
  var _normalizeForms = ["NFD", "NFC", "NFKD", "NFKC"].reduce((accum, form) => {
    try {
      if ("test".normalize(form) !== "test") {
        throw new Error("bad");
      }
      ;
      if (form === "NFD") {
        const check = String.fromCharCode(233).normalize("NFD");
        const expected = String.fromCharCode(101, 769);
        if (check !== expected) {
          throw new Error("broken");
        }
      }
      accum.push(form);
    } catch (error) {
    }
    return accum;
  }, []);
  function assertNormalize(form) {
    assert(_normalizeForms.indexOf(form) >= 0, "platform missing String.prototype.normalize", "UNSUPPORTED_OPERATION", {
      operation: "String.prototype.normalize",
      info: { form }
    });
  }
  function assertPrivate(givenGuard, guard, className) {
    if (className == null) {
      className = "";
    }
    if (givenGuard !== guard) {
      let method = className, operation = "new";
      if (className) {
        method += ".";
        operation += " " + className;
      }
      assert(false, `private constructor; use ${method}from* methods`, "UNSUPPORTED_OPERATION", {
        operation
      });
    }
  }

  // node_modules/ethers/lib.esm/utils/data.js
  function _getBytes(value, name, copy4) {
    if (value instanceof Uint8Array) {
      if (copy4) {
        return new Uint8Array(value);
      }
      return value;
    }
    if (typeof value === "string" && value.match(/^0x([0-9a-f][0-9a-f])*$/i)) {
      const result = new Uint8Array((value.length - 2) / 2);
      let offset = 2;
      for (let i4 = 0; i4 < result.length; i4++) {
        result[i4] = parseInt(value.substring(offset, offset + 2), 16);
        offset += 2;
      }
      return result;
    }
    assertArgument(false, "invalid BytesLike value", name || "value", value);
  }
  function getBytes(value, name) {
    return _getBytes(value, name, false);
  }
  function getBytesCopy(value, name) {
    return _getBytes(value, name, true);
  }
  function isHexString(value, length) {
    if (typeof value !== "string" || !value.match(/^0x[0-9A-Fa-f]*$/)) {
      return false;
    }
    if (typeof length === "number" && value.length !== 2 + 2 * length) {
      return false;
    }
    if (length === true && value.length % 2 !== 0) {
      return false;
    }
    return true;
  }
  function isBytesLike(value) {
    return isHexString(value, true) || value instanceof Uint8Array;
  }
  var HexCharacters = "0123456789abcdef";
  function hexlify(data) {
    const bytes2 = getBytes(data);
    let result = "0x";
    for (let i4 = 0; i4 < bytes2.length; i4++) {
      const v3 = bytes2[i4];
      result += HexCharacters[(v3 & 240) >> 4] + HexCharacters[v3 & 15];
    }
    return result;
  }
  function concat(datas) {
    return "0x" + datas.map((d4) => hexlify(d4).substring(2)).join("");
  }
  function dataLength(data) {
    if (isHexString(data, true)) {
      return (data.length - 2) / 2;
    }
    return getBytes(data).length;
  }
  function dataSlice(data, start, end) {
    const bytes2 = getBytes(data);
    if (end != null && end > bytes2.length) {
      assert(false, "cannot slice beyond data bounds", "BUFFER_OVERRUN", {
        buffer: bytes2,
        length: bytes2.length,
        offset: end
      });
    }
    return hexlify(bytes2.slice(start == null ? 0 : start, end == null ? bytes2.length : end));
  }
  function zeroPad(data, length, left) {
    const bytes2 = getBytes(data);
    assert(length >= bytes2.length, "padding exceeds data length", "BUFFER_OVERRUN", {
      buffer: new Uint8Array(bytes2),
      length,
      offset: length + 1
    });
    const result = new Uint8Array(length);
    result.fill(0);
    if (left) {
      result.set(bytes2, length - bytes2.length);
    } else {
      result.set(bytes2, 0);
    }
    return hexlify(result);
  }
  function zeroPadValue(data, length) {
    return zeroPad(data, length, true);
  }
  function zeroPadBytes(data, length) {
    return zeroPad(data, length, false);
  }

  // node_modules/ethers/lib.esm/utils/maths.js
  var BN_0 = BigInt(0);
  var BN_1 = BigInt(1);
  var maxValue = 9007199254740991;
  function fromTwos(_value, _width) {
    const value = getUint(_value, "value");
    const width = BigInt(getNumber(_width, "width"));
    assert(value >> width === BN_0, "overflow", "NUMERIC_FAULT", {
      operation: "fromTwos",
      fault: "overflow",
      value: _value
    });
    if (value >> width - BN_1) {
      const mask2 = (BN_1 << width) - BN_1;
      return -((~value & mask2) + BN_1);
    }
    return value;
  }
  function toTwos(_value, _width) {
    let value = getBigInt(_value, "value");
    const width = BigInt(getNumber(_width, "width"));
    const limit = BN_1 << width - BN_1;
    if (value < BN_0) {
      value = -value;
      assert(value <= limit, "too low", "NUMERIC_FAULT", {
        operation: "toTwos",
        fault: "overflow",
        value: _value
      });
      const mask2 = (BN_1 << width) - BN_1;
      return (~value & mask2) + BN_1;
    } else {
      assert(value < limit, "too high", "NUMERIC_FAULT", {
        operation: "toTwos",
        fault: "overflow",
        value: _value
      });
    }
    return value;
  }
  function mask(_value, _bits) {
    const value = getUint(_value, "value");
    const bits = BigInt(getNumber(_bits, "bits"));
    return value & (BN_1 << bits) - BN_1;
  }
  function getBigInt(value, name) {
    switch (typeof value) {
      case "bigint":
        return value;
      case "number":
        assertArgument(Number.isInteger(value), "underflow", name || "value", value);
        assertArgument(value >= -maxValue && value <= maxValue, "overflow", name || "value", value);
        return BigInt(value);
      case "string":
        try {
          if (value === "") {
            throw new Error("empty string");
          }
          if (value[0] === "-" && value[1] !== "-") {
            return -BigInt(value.substring(1));
          }
          return BigInt(value);
        } catch (e4) {
          assertArgument(false, `invalid BigNumberish string: ${e4.message}`, name || "value", value);
        }
    }
    assertArgument(false, "invalid BigNumberish value", name || "value", value);
  }
  function getUint(value, name) {
    const result = getBigInt(value, name);
    assert(result >= BN_0, "unsigned value cannot be negative", "NUMERIC_FAULT", {
      fault: "overflow",
      operation: "getUint",
      value
    });
    return result;
  }
  var Nibbles = "0123456789abcdef";
  function toBigInt(value) {
    if (value instanceof Uint8Array) {
      let result = "0x0";
      for (const v3 of value) {
        result += Nibbles[v3 >> 4];
        result += Nibbles[v3 & 15];
      }
      return BigInt(result);
    }
    return getBigInt(value);
  }
  function getNumber(value, name) {
    switch (typeof value) {
      case "bigint":
        assertArgument(value >= -maxValue && value <= maxValue, "overflow", name || "value", value);
        return Number(value);
      case "number":
        assertArgument(Number.isInteger(value), "underflow", name || "value", value);
        assertArgument(value >= -maxValue && value <= maxValue, "overflow", name || "value", value);
        return value;
      case "string":
        try {
          if (value === "") {
            throw new Error("empty string");
          }
          return getNumber(BigInt(value), name);
        } catch (e4) {
          assertArgument(false, `invalid numeric string: ${e4.message}`, name || "value", value);
        }
    }
    assertArgument(false, "invalid numeric value", name || "value", value);
  }
  function toNumber(value) {
    return getNumber(toBigInt(value));
  }
  function toBeHex(_value, _width) {
    const value = getUint(_value, "value");
    let result = value.toString(16);
    if (_width == null) {
      if (result.length % 2) {
        result = "0" + result;
      }
    } else {
      const width = getNumber(_width, "width");
      assert(width * 2 >= result.length, `value exceeds width (${width} bits)`, "NUMERIC_FAULT", {
        operation: "toBeHex",
        fault: "overflow",
        value: _value
      });
      while (result.length < width * 2) {
        result = "0" + result;
      }
    }
    return "0x" + result;
  }
  function toBeArray(_value) {
    const value = getUint(_value, "value");
    if (value === BN_0) {
      return new Uint8Array([]);
    }
    let hex = value.toString(16);
    if (hex.length % 2) {
      hex = "0" + hex;
    }
    const result = new Uint8Array(hex.length / 2);
    for (let i4 = 0; i4 < result.length; i4++) {
      const offset = i4 * 2;
      result[i4] = parseInt(hex.substring(offset, offset + 2), 16);
    }
    return result;
  }
  function toQuantity(value) {
    let result = hexlify(isBytesLike(value) ? value : toBeArray(value)).substring(2);
    while (result.startsWith("0")) {
      result = result.substring(1);
    }
    if (result === "") {
      result = "0";
    }
    return "0x" + result;
  }

  // node_modules/ethers/lib.esm/utils/base58.js
  var Alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz";
  var BN_02 = BigInt(0);
  var BN_58 = BigInt(58);
  function encodeBase58(_value) {
    let value = toBigInt(getBytes(_value));
    let result = "";
    while (value) {
      result = Alphabet[Number(value % BN_58)] + result;
      value /= BN_58;
    }
    return result;
  }

  // node_modules/ethers/lib.esm/utils/base64-browser.js
  function decodeBase64(textData) {
    textData = atob(textData);
    const data = new Uint8Array(textData.length);
    for (let i4 = 0; i4 < textData.length; i4++) {
      data[i4] = textData.charCodeAt(i4);
    }
    return getBytes(data);
  }
  function encodeBase64(_data) {
    const data = getBytes(_data);
    let textData = "";
    for (let i4 = 0; i4 < data.length; i4++) {
      textData += String.fromCharCode(data[i4]);
    }
    return btoa(textData);
  }

  // node_modules/ethers/lib.esm/utils/events.js
  var EventPayload = class {
    /**
     *  The event filter.
     */
    filter;
    /**
     *  The **EventEmitterable**.
     */
    emitter;
    #listener;
    /**
     *  Create a new **EventPayload** for %%emitter%% with
     *  the %%listener%% and for %%filter%%.
     */
    constructor(emitter, listener, filter) {
      this.#listener = listener;
      defineProperties(this, { emitter, filter });
    }
    /**
     *  Unregister the triggered listener for future events.
     */
    async removeListener() {
      if (this.#listener == null) {
        return;
      }
      await this.emitter.off(this.filter, this.#listener);
    }
  };

  // node_modules/ethers/lib.esm/utils/utf8.js
  function errorFunc(reason, offset, bytes2, output2, badCodepoint) {
    assertArgument(false, `invalid codepoint at offset ${offset}; ${reason}`, "bytes", bytes2);
  }
  function ignoreFunc(reason, offset, bytes2, output2, badCodepoint) {
    if (reason === "BAD_PREFIX" || reason === "UNEXPECTED_CONTINUE") {
      let i4 = 0;
      for (let o4 = offset + 1; o4 < bytes2.length; o4++) {
        if (bytes2[o4] >> 6 !== 2) {
          break;
        }
        i4++;
      }
      return i4;
    }
    if (reason === "OVERRUN") {
      return bytes2.length - offset - 1;
    }
    return 0;
  }
  function replaceFunc(reason, offset, bytes2, output2, badCodepoint) {
    if (reason === "OVERLONG") {
      assertArgument(typeof badCodepoint === "number", "invalid bad code point for replacement", "badCodepoint", badCodepoint);
      output2.push(badCodepoint);
      return 0;
    }
    output2.push(65533);
    return ignoreFunc(reason, offset, bytes2, output2, badCodepoint);
  }
  var Utf8ErrorFuncs = Object.freeze({
    error: errorFunc,
    ignore: ignoreFunc,
    replace: replaceFunc
  });
  function getUtf8CodePoints(_bytes, onError) {
    if (onError == null) {
      onError = Utf8ErrorFuncs.error;
    }
    const bytes2 = getBytes(_bytes, "bytes");
    const result = [];
    let i4 = 0;
    while (i4 < bytes2.length) {
      const c4 = bytes2[i4++];
      if (c4 >> 7 === 0) {
        result.push(c4);
        continue;
      }
      let extraLength = null;
      let overlongMask = null;
      if ((c4 & 224) === 192) {
        extraLength = 1;
        overlongMask = 127;
      } else if ((c4 & 240) === 224) {
        extraLength = 2;
        overlongMask = 2047;
      } else if ((c4 & 248) === 240) {
        extraLength = 3;
        overlongMask = 65535;
      } else {
        if ((c4 & 192) === 128) {
          i4 += onError("UNEXPECTED_CONTINUE", i4 - 1, bytes2, result);
        } else {
          i4 += onError("BAD_PREFIX", i4 - 1, bytes2, result);
        }
        continue;
      }
      if (i4 - 1 + extraLength >= bytes2.length) {
        i4 += onError("OVERRUN", i4 - 1, bytes2, result);
        continue;
      }
      let res = c4 & (1 << 8 - extraLength - 1) - 1;
      for (let j4 = 0; j4 < extraLength; j4++) {
        let nextChar = bytes2[i4];
        if ((nextChar & 192) != 128) {
          i4 += onError("MISSING_CONTINUE", i4, bytes2, result);
          res = null;
          break;
        }
        ;
        res = res << 6 | nextChar & 63;
        i4++;
      }
      if (res === null) {
        continue;
      }
      if (res > 1114111) {
        i4 += onError("OUT_OF_RANGE", i4 - 1 - extraLength, bytes2, result, res);
        continue;
      }
      if (res >= 55296 && res <= 57343) {
        i4 += onError("UTF16_SURROGATE", i4 - 1 - extraLength, bytes2, result, res);
        continue;
      }
      if (res <= overlongMask) {
        i4 += onError("OVERLONG", i4 - 1 - extraLength, bytes2, result, res);
        continue;
      }
      result.push(res);
    }
    return result;
  }
  function toUtf8Bytes(str, form) {
    if (form != null) {
      assertNormalize(form);
      str = str.normalize(form);
    }
    let result = [];
    for (let i4 = 0; i4 < str.length; i4++) {
      const c4 = str.charCodeAt(i4);
      if (c4 < 128) {
        result.push(c4);
      } else if (c4 < 2048) {
        result.push(c4 >> 6 | 192);
        result.push(c4 & 63 | 128);
      } else if ((c4 & 64512) == 55296) {
        i4++;
        const c22 = str.charCodeAt(i4);
        assertArgument(i4 < str.length && (c22 & 64512) === 56320, "invalid surrogate pair", "str", str);
        const pair = 65536 + ((c4 & 1023) << 10) + (c22 & 1023);
        result.push(pair >> 18 | 240);
        result.push(pair >> 12 & 63 | 128);
        result.push(pair >> 6 & 63 | 128);
        result.push(pair & 63 | 128);
      } else {
        result.push(c4 >> 12 | 224);
        result.push(c4 >> 6 & 63 | 128);
        result.push(c4 & 63 | 128);
      }
    }
    return new Uint8Array(result);
  }
  function _toUtf8String(codePoints) {
    return codePoints.map((codePoint) => {
      if (codePoint <= 65535) {
        return String.fromCharCode(codePoint);
      }
      codePoint -= 65536;
      return String.fromCharCode((codePoint >> 10 & 1023) + 55296, (codePoint & 1023) + 56320);
    }).join("");
  }
  function toUtf8String(bytes2, onError) {
    return _toUtf8String(getUtf8CodePoints(bytes2, onError));
  }

  // node_modules/ethers/lib.esm/utils/geturl-browser.js
  async function getUrl(req, _signal) {
    const protocol = req.url.split(":")[0].toLowerCase();
    assert(protocol === "http" || protocol === "https", `unsupported protocol ${protocol}`, "UNSUPPORTED_OPERATION", {
      info: { protocol },
      operation: "request"
    });
    assert(protocol === "https" || !req.credentials || req.allowInsecureAuthentication, "insecure authorized connections unsupported", "UNSUPPORTED_OPERATION", {
      operation: "request"
    });
    let signal = void 0;
    if (_signal) {
      const controller = new AbortController();
      signal = controller.signal;
      _signal.addListener(() => {
        controller.abort();
      });
    }
    const init = {
      method: req.method,
      headers: new Headers(Array.from(req)),
      body: req.body || void 0,
      signal
    };
    const resp = await fetch(req.url, init);
    const headers = {};
    resp.headers.forEach((value, key) => {
      headers[key.toLowerCase()] = value;
    });
    const respBody = await resp.arrayBuffer();
    const body = respBody == null ? null : new Uint8Array(respBody);
    return {
      statusCode: resp.status,
      statusMessage: resp.statusText,
      headers,
      body
    };
  }

  // node_modules/ethers/lib.esm/utils/fetch.js
  var MAX_ATTEMPTS = 12;
  var SLOT_INTERVAL = 250;
  var getUrlFunc = getUrl;
  var reData = new RegExp("^data:([^;:]*)?(;base64)?,(.*)$", "i");
  var reIpfs = new RegExp("^ipfs://(ipfs/)?(.*)$", "i");
  var locked = false;
  async function dataGatewayFunc(url, signal) {
    try {
      const match = url.match(reData);
      if (!match) {
        throw new Error("invalid data");
      }
      return new FetchResponse(200, "OK", {
        "content-type": match[1] || "text/plain"
      }, match[2] ? decodeBase64(match[3]) : unpercent(match[3]));
    } catch (error) {
      return new FetchResponse(599, "BAD REQUEST (invalid data: URI)", {}, null, new FetchRequest(url));
    }
  }
  function getIpfsGatewayFunc(baseUrl) {
    async function gatewayIpfs(url, signal) {
      try {
        const match = url.match(reIpfs);
        if (!match) {
          throw new Error("invalid link");
        }
        return new FetchRequest(`${baseUrl}${match[2]}`);
      } catch (error) {
        return new FetchResponse(599, "BAD REQUEST (invalid IPFS URI)", {}, null, new FetchRequest(url));
      }
    }
    return gatewayIpfs;
  }
  var Gateways = {
    "data": dataGatewayFunc,
    "ipfs": getIpfsGatewayFunc("https://gateway.ipfs.io/ipfs/")
  };
  var fetchSignals = /* @__PURE__ */ new WeakMap();
  var FetchCancelSignal = class {
    #listeners;
    #cancelled;
    constructor(request) {
      this.#listeners = [];
      this.#cancelled = false;
      fetchSignals.set(request, () => {
        if (this.#cancelled) {
          return;
        }
        this.#cancelled = true;
        for (const listener of this.#listeners) {
          setTimeout(() => {
            listener();
          }, 0);
        }
        this.#listeners = [];
      });
    }
    addListener(listener) {
      assert(!this.#cancelled, "singal already cancelled", "UNSUPPORTED_OPERATION", {
        operation: "fetchCancelSignal.addCancelListener"
      });
      this.#listeners.push(listener);
    }
    get cancelled() {
      return this.#cancelled;
    }
    checkSignal() {
      assert(!this.cancelled, "cancelled", "CANCELLED", {});
    }
  };
  function checkSignal(signal) {
    if (signal == null) {
      throw new Error("missing signal; should not happen");
    }
    signal.checkSignal();
    return signal;
  }
  var FetchRequest = class _FetchRequest {
    #allowInsecure;
    #gzip;
    #headers;
    #method;
    #timeout;
    #url;
    #body;
    #bodyType;
    #creds;
    // Hooks
    #preflight;
    #process;
    #retry;
    #signal;
    #throttle;
    /**
     *  The fetch URI to requrest.
     */
    get url() {
      return this.#url;
    }
    set url(url) {
      this.#url = String(url);
    }
    /**
     *  The fetch body, if any, to send as the request body. //(default: null)//
     *
     *  When setting a body, the intrinsic ``Content-Type`` is automatically
     *  set and will be used if **not overridden** by setting a custom
     *  header.
     *
     *  If %%body%% is null, the body is cleared (along with the
     *  intrinsic ``Content-Type``) and the .
     *
     *  If %%body%% is a string, the intrincis ``Content-Type`` is set to
     *  ``text/plain``.
     *
     *  If %%body%% is a Uint8Array, the intrincis ``Content-Type`` is set to
     *  ``application/octet-stream``.
     *
     *  If %%body%% is any other object, the intrincis ``Content-Type`` is
     *  set to ``application/json``.
     */
    get body() {
      if (this.#body == null) {
        return null;
      }
      return new Uint8Array(this.#body);
    }
    set body(body) {
      if (body == null) {
        this.#body = void 0;
        this.#bodyType = void 0;
      } else if (typeof body === "string") {
        this.#body = toUtf8Bytes(body);
        this.#bodyType = "text/plain";
      } else if (body instanceof Uint8Array) {
        this.#body = body;
        this.#bodyType = "application/octet-stream";
      } else if (typeof body === "object") {
        this.#body = toUtf8Bytes(JSON.stringify(body));
        this.#bodyType = "application/json";
      } else {
        throw new Error("invalid body");
      }
    }
    /**
     *  Returns true if the request has a body.
     */
    hasBody() {
      return this.#body != null;
    }
    /**
     *  The HTTP method to use when requesting the URI. If no method
     *  has been explicitly set, then ``GET`` is used if the body is
     *  null and ``POST`` otherwise.
     */
    get method() {
      if (this.#method) {
        return this.#method;
      }
      if (this.hasBody()) {
        return "POST";
      }
      return "GET";
    }
    set method(method) {
      if (method == null) {
        method = "";
      }
      this.#method = String(method).toUpperCase();
    }
    /**
     *  The headers that will be used when requesting the URI. All
     *  keys are lower-case.
     *
     *  This object is a copy, so any chnages will **NOT** be reflected
     *  in the ``FetchRequest``.
     *
     *  To set a header entry, use the ``setHeader`` method.
     */
    get headers() {
      const headers = Object.assign({}, this.#headers);
      if (this.#creds) {
        headers["authorization"] = `Basic ${encodeBase64(toUtf8Bytes(this.#creds))}`;
      }
      ;
      if (this.allowGzip) {
        headers["accept-encoding"] = "gzip";
      }
      if (headers["content-type"] == null && this.#bodyType) {
        headers["content-type"] = this.#bodyType;
      }
      if (this.body) {
        headers["content-length"] = String(this.body.length);
      }
      return headers;
    }
    /**
     *  Get the header for %%key%%, ignoring case.
     */
    getHeader(key) {
      return this.headers[key.toLowerCase()];
    }
    /**
     *  Set the header for %%key%% to %%value%%. All values are coerced
     *  to a string.
     */
    setHeader(key, value) {
      this.#headers[String(key).toLowerCase()] = String(value);
    }
    /**
     *  Clear all headers, resetting all intrinsic headers.
     */
    clearHeaders() {
      this.#headers = {};
    }
    [Symbol.iterator]() {
      const headers = this.headers;
      const keys = Object.keys(headers);
      let index = 0;
      return {
        next: () => {
          if (index < keys.length) {
            const key = keys[index++];
            return {
              value: [key, headers[key]],
              done: false
            };
          }
          return { value: void 0, done: true };
        }
      };
    }
    /**
     *  The value that will be sent for the ``Authorization`` header.
     *
     *  To set the credentials, use the ``setCredentials`` method.
     */
    get credentials() {
      return this.#creds || null;
    }
    /**
     *  Sets an ``Authorization`` for %%username%% with %%password%%.
     */
    setCredentials(username, password) {
      assertArgument(!username.match(/:/), "invalid basic authentication username", "username", "[REDACTED]");
      this.#creds = `${username}:${password}`;
    }
    /**
     *  Enable and request gzip-encoded responses. The response will
     *  automatically be decompressed. //(default: true)//
     */
    get allowGzip() {
      return this.#gzip;
    }
    set allowGzip(value) {
      this.#gzip = !!value;
    }
    /**
     *  Allow ``Authentication`` credentials to be sent over insecure
     *  channels. //(default: false)//
     */
    get allowInsecureAuthentication() {
      return !!this.#allowInsecure;
    }
    set allowInsecureAuthentication(value) {
      this.#allowInsecure = !!value;
    }
    /**
     *  The timeout (in milliseconds) to wait for a complere response.
     *  //(default: 5 minutes)//
     */
    get timeout() {
      return this.#timeout;
    }
    set timeout(timeout) {
      assertArgument(timeout >= 0, "timeout must be non-zero", "timeout", timeout);
      this.#timeout = timeout;
    }
    /**
     *  This function is called prior to each request, for example
     *  during a redirection or retry in case of server throttling.
     *
     *  This offers an opportunity to populate headers or update
     *  content before sending a request.
     */
    get preflightFunc() {
      return this.#preflight || null;
    }
    set preflightFunc(preflight) {
      this.#preflight = preflight;
    }
    /**
     *  This function is called after each response, offering an
     *  opportunity to provide client-level throttling or updating
     *  response data.
     *
     *  Any error thrown in this causes the ``send()`` to throw.
     *
     *  To schedule a retry attempt (assuming the maximum retry limit
     *  has not been reached), use [[response.throwThrottleError]].
     */
    get processFunc() {
      return this.#process || null;
    }
    set processFunc(process2) {
      this.#process = process2;
    }
    /**
     *  This function is called on each retry attempt.
     */
    get retryFunc() {
      return this.#retry || null;
    }
    set retryFunc(retry) {
      this.#retry = retry;
    }
    /**
     *  Create a new FetchRequest instance with default values.
     *
     *  Once created, each property may be set before issuing a
     *  ``.send()`` to make the request.
     */
    constructor(url) {
      this.#url = String(url);
      this.#allowInsecure = false;
      this.#gzip = true;
      this.#headers = {};
      this.#method = "";
      this.#timeout = 3e5;
      this.#throttle = {
        slotInterval: SLOT_INTERVAL,
        maxAttempts: MAX_ATTEMPTS
      };
    }
    toString() {
      return `<FetchRequest method=${JSON.stringify(this.method)} url=${JSON.stringify(this.url)} headers=${JSON.stringify(this.headers)} body=${this.#body ? hexlify(this.#body) : "null"}>`;
    }
    /**
     *  Update the throttle parameters used to determine maximum
     *  attempts and exponential-backoff properties.
     */
    setThrottleParams(params) {
      if (params.slotInterval != null) {
        this.#throttle.slotInterval = params.slotInterval;
      }
      if (params.maxAttempts != null) {
        this.#throttle.maxAttempts = params.maxAttempts;
      }
    }
    async #send(attempt, expires, delay, _request, _response) {
      if (attempt >= this.#throttle.maxAttempts) {
        return _response.makeServerError("exceeded maximum retry limit");
      }
      assert(getTime() <= expires, "timeout", "TIMEOUT", {
        operation: "request.send",
        reason: "timeout",
        request: _request
      });
      if (delay > 0) {
        await wait(delay);
      }
      let req = this.clone();
      const scheme = (req.url.split(":")[0] || "").toLowerCase();
      if (scheme in Gateways) {
        const result = await Gateways[scheme](req.url, checkSignal(_request.#signal));
        if (result instanceof FetchResponse) {
          let response2 = result;
          if (this.processFunc) {
            checkSignal(_request.#signal);
            try {
              response2 = await this.processFunc(req, response2);
            } catch (error) {
              if (error.throttle == null || typeof error.stall !== "number") {
                response2.makeServerError("error in post-processing function", error).assertOk();
              }
            }
          }
          return response2;
        }
        req = result;
      }
      if (this.preflightFunc) {
        req = await this.preflightFunc(req);
      }
      const resp = await getUrlFunc(req, checkSignal(_request.#signal));
      let response = new FetchResponse(resp.statusCode, resp.statusMessage, resp.headers, resp.body, _request);
      if (response.statusCode === 301 || response.statusCode === 302) {
        try {
          const location2 = response.headers.location || "";
          return req.redirect(location2).#send(attempt + 1, expires, 0, _request, response);
        } catch (error) {
        }
        return response;
      } else if (response.statusCode === 429) {
        if (this.retryFunc == null || await this.retryFunc(req, response, attempt)) {
          const retryAfter = response.headers["retry-after"];
          let delay2 = this.#throttle.slotInterval * Math.trunc(Math.random() * Math.pow(2, attempt));
          if (typeof retryAfter === "string" && retryAfter.match(/^[1-9][0-9]*$/)) {
            delay2 = parseInt(retryAfter);
          }
          return req.clone().#send(attempt + 1, expires, delay2, _request, response);
        }
      }
      if (this.processFunc) {
        checkSignal(_request.#signal);
        try {
          response = await this.processFunc(req, response);
        } catch (error) {
          if (error.throttle == null || typeof error.stall !== "number") {
            response.makeServerError("error in post-processing function", error).assertOk();
          }
          let delay2 = this.#throttle.slotInterval * Math.trunc(Math.random() * Math.pow(2, attempt));
          ;
          if (error.stall >= 0) {
            delay2 = error.stall;
          }
          return req.clone().#send(attempt + 1, expires, delay2, _request, response);
        }
      }
      return response;
    }
    /**
     *  Resolves to the response by sending the request.
     */
    send() {
      assert(this.#signal == null, "request already sent", "UNSUPPORTED_OPERATION", { operation: "fetchRequest.send" });
      this.#signal = new FetchCancelSignal(this);
      return this.#send(0, getTime() + this.timeout, 0, this, new FetchResponse(0, "", {}, null, this));
    }
    /**
     *  Cancels the inflight response, causing a ``CANCELLED``
     *  error to be rejected from the [[send]].
     */
    cancel() {
      assert(this.#signal != null, "request has not been sent", "UNSUPPORTED_OPERATION", { operation: "fetchRequest.cancel" });
      const signal = fetchSignals.get(this);
      if (!signal) {
        throw new Error("missing signal; should not happen");
      }
      signal();
    }
    /**
     *  Returns a new [[FetchRequest]] that represents the redirection
     *  to %%location%%.
     */
    redirect(location2) {
      const current = this.url.split(":")[0].toLowerCase();
      const target = location2.split(":")[0].toLowerCase();
      assert(this.method === "GET" && (current !== "https" || target !== "http") && location2.match(/^https?:/), `unsupported redirect`, "UNSUPPORTED_OPERATION", {
        operation: `redirect(${this.method} ${JSON.stringify(this.url)} => ${JSON.stringify(location2)})`
      });
      const req = new _FetchRequest(location2);
      req.method = "GET";
      req.allowGzip = this.allowGzip;
      req.timeout = this.timeout;
      req.#headers = Object.assign({}, this.#headers);
      if (this.#body) {
        req.#body = new Uint8Array(this.#body);
      }
      req.#bodyType = this.#bodyType;
      return req;
    }
    /**
     *  Create a new copy of this request.
     */
    clone() {
      const clone = new _FetchRequest(this.url);
      clone.#method = this.#method;
      if (this.#body) {
        clone.#body = this.#body;
      }
      clone.#bodyType = this.#bodyType;
      clone.#headers = Object.assign({}, this.#headers);
      clone.#creds = this.#creds;
      if (this.allowGzip) {
        clone.allowGzip = true;
      }
      clone.timeout = this.timeout;
      if (this.allowInsecureAuthentication) {
        clone.allowInsecureAuthentication = true;
      }
      clone.#preflight = this.#preflight;
      clone.#process = this.#process;
      clone.#retry = this.#retry;
      return clone;
    }
    /**
     *  Locks all static configuration for gateways and FetchGetUrlFunc
     *  registration.
     */
    static lockConfig() {
      locked = true;
    }
    /**
     *  Get the current Gateway function for %%scheme%%.
     */
    static getGateway(scheme) {
      return Gateways[scheme.toLowerCase()] || null;
    }
    /**
     *  Use the %%func%% when fetching URIs using %%scheme%%.
     *
     *  This method affects all requests globally.
     *
     *  If [[lockConfig]] has been called, no change is made and this
     *  throws.
     */
    static registerGateway(scheme, func) {
      scheme = scheme.toLowerCase();
      if (scheme === "http" || scheme === "https") {
        throw new Error(`cannot intercept ${scheme}; use registerGetUrl`);
      }
      if (locked) {
        throw new Error("gateways locked");
      }
      Gateways[scheme] = func;
    }
    /**
     *  Use %%getUrl%% when fetching URIs over HTTP and HTTPS requests.
     *
     *  This method affects all requests globally.
     *
     *  If [[lockConfig]] has been called, no change is made and this
     *  throws.
     */
    static registerGetUrl(getUrl2) {
      if (locked) {
        throw new Error("gateways locked");
      }
      getUrlFunc = getUrl2;
    }
    /**
     *  Creates a function that can "fetch" data URIs.
     *
     *  Note that this is automatically done internally to support
     *  data URIs, so it is not necessary to register it.
     *
     *  This is not generally something that is needed, but may
     *  be useful in a wrapper to perfom custom data URI functionality.
     */
    static createDataGateway() {
      return dataGatewayFunc;
    }
    /**
     *  Creates a function that will fetch IPFS (unvalidated) from
     *  a custom gateway baseUrl.
     *
     *  The default IPFS gateway used internally is
     *  ``"https:/\/gateway.ipfs.io/ipfs/"``.
     */
    static createIpfsGatewayFunc(baseUrl) {
      return getIpfsGatewayFunc(baseUrl);
    }
  };
  var FetchResponse = class _FetchResponse {
    #statusCode;
    #statusMessage;
    #headers;
    #body;
    #request;
    #error;
    toString() {
      return `<FetchResponse status=${this.statusCode} body=${this.#body ? hexlify(this.#body) : "null"}>`;
    }
    /**
     *  The response status code.
     */
    get statusCode() {
      return this.#statusCode;
    }
    /**
     *  The response status message.
     */
    get statusMessage() {
      return this.#statusMessage;
    }
    /**
     *  The response headers. All keys are lower-case.
     */
    get headers() {
      return Object.assign({}, this.#headers);
    }
    /**
     *  The response body, or ``null`` if there was no body.
     */
    get body() {
      return this.#body == null ? null : new Uint8Array(this.#body);
    }
    /**
     *  The response body as a UTF-8 encoded string, or the empty
     *  string (i.e. ``""``) if there was no body.
     *
     *  An error is thrown if the body is invalid UTF-8 data.
     */
    get bodyText() {
      try {
        return this.#body == null ? "" : toUtf8String(this.#body);
      } catch (error) {
        assert(false, "response body is not valid UTF-8 data", "UNSUPPORTED_OPERATION", {
          operation: "bodyText",
          info: { response: this }
        });
      }
    }
    /**
     *  The response body, decoded as JSON.
     *
     *  An error is thrown if the body is invalid JSON-encoded data
     *  or if there was no body.
     */
    get bodyJson() {
      try {
        return JSON.parse(this.bodyText);
      } catch (error) {
        assert(false, "response body is not valid JSON", "UNSUPPORTED_OPERATION", {
          operation: "bodyJson",
          info: { response: this }
        });
      }
    }
    [Symbol.iterator]() {
      const headers = this.headers;
      const keys = Object.keys(headers);
      let index = 0;
      return {
        next: () => {
          if (index < keys.length) {
            const key = keys[index++];
            return {
              value: [key, headers[key]],
              done: false
            };
          }
          return { value: void 0, done: true };
        }
      };
    }
    constructor(statusCode, statusMessage, headers, body, request) {
      this.#statusCode = statusCode;
      this.#statusMessage = statusMessage;
      this.#headers = Object.keys(headers).reduce((accum, k3) => {
        accum[k3.toLowerCase()] = String(headers[k3]);
        return accum;
      }, {});
      this.#body = body == null ? null : new Uint8Array(body);
      this.#request = request || null;
      this.#error = { message: "" };
    }
    /**
     *  Return a Response with matching headers and body, but with
     *  an error status code (i.e. 599) and %%message%% with an
     *  optional %%error%%.
     */
    makeServerError(message, error) {
      let statusMessage;
      if (!message) {
        message = `${this.statusCode} ${this.statusMessage}`;
        statusMessage = `CLIENT ESCALATED SERVER ERROR (${message})`;
      } else {
        statusMessage = `CLIENT ESCALATED SERVER ERROR (${this.statusCode} ${this.statusMessage}; ${message})`;
      }
      const response = new _FetchResponse(599, statusMessage, this.headers, this.body, this.#request || void 0);
      response.#error = { message, error };
      return response;
    }
    /**
     *  If called within a [request.processFunc](FetchRequest-processFunc)
     *  call, causes the request to retry as if throttled for %%stall%%
     *  milliseconds.
     */
    throwThrottleError(message, stall2) {
      if (stall2 == null) {
        stall2 = -1;
      } else {
        assertArgument(Number.isInteger(stall2) && stall2 >= 0, "invalid stall timeout", "stall", stall2);
      }
      const error = new Error(message || "throttling requests");
      defineProperties(error, { stall: stall2, throttle: true });
      throw error;
    }
    /**
     *  Get the header value for %%key%%, ignoring case.
     */
    getHeader(key) {
      return this.headers[key.toLowerCase()];
    }
    /**
     *  Returns true of the response has a body.
     */
    hasBody() {
      return this.#body != null;
    }
    /**
     *  The request made for this response.
     */
    get request() {
      return this.#request;
    }
    /**
     *  Returns true if this response was a success statusCode.
     */
    ok() {
      return this.#error.message === "" && this.statusCode >= 200 && this.statusCode < 300;
    }
    /**
     *  Throws a ``SERVER_ERROR`` if this response is not ok.
     */
    assertOk() {
      if (this.ok()) {
        return;
      }
      let { message, error } = this.#error;
      if (message === "") {
        message = `server response ${this.statusCode} ${this.statusMessage}`;
      }
      assert(false, message, "SERVER_ERROR", {
        request: this.request || "unknown request",
        response: this,
        error
      });
    }
  };
  function getTime() {
    return (/* @__PURE__ */ new Date()).getTime();
  }
  function unpercent(value) {
    return toUtf8Bytes(value.replace(/%([0-9a-f][0-9a-f])/gi, (all, code) => {
      return String.fromCharCode(parseInt(code, 16));
    }));
  }
  function wait(delay) {
    return new Promise((resolve) => setTimeout(resolve, delay));
  }

  // node_modules/ethers/lib.esm/utils/fixednumber.js
  var BN_N1 = BigInt(-1);
  var BN_03 = BigInt(0);
  var BN_12 = BigInt(1);
  var BN_5 = BigInt(5);
  var _guard = {};
  var Zeros = "0000";
  while (Zeros.length < 80) {
    Zeros += Zeros;
  }
  function getTens(decimals) {
    let result = Zeros;
    while (result.length < decimals) {
      result += result;
    }
    return BigInt("1" + result.substring(0, decimals));
  }
  function checkValue(val, format, safeOp) {
    const width = BigInt(format.width);
    if (format.signed) {
      const limit = BN_12 << width - BN_12;
      assert(safeOp == null || val >= -limit && val < limit, "overflow", "NUMERIC_FAULT", {
        operation: safeOp,
        fault: "overflow",
        value: val
      });
      if (val > BN_03) {
        val = fromTwos(mask(val, width), width);
      } else {
        val = -fromTwos(mask(-val, width), width);
      }
    } else {
      const limit = BN_12 << width;
      assert(safeOp == null || val >= 0 && val < limit, "overflow", "NUMERIC_FAULT", {
        operation: safeOp,
        fault: "overflow",
        value: val
      });
      val = (val % limit + limit) % limit & limit - BN_12;
    }
    return val;
  }
  function getFormat(value) {
    if (typeof value === "number") {
      value = `fixed128x${value}`;
    }
    let signed2 = true;
    let width = 128;
    let decimals = 18;
    if (typeof value === "string") {
      if (value === "fixed") {
      } else if (value === "ufixed") {
        signed2 = false;
      } else {
        const match = value.match(/^(u?)fixed([0-9]+)x([0-9]+)$/);
        assertArgument(match, "invalid fixed format", "format", value);
        signed2 = match[1] !== "u";
        width = parseInt(match[2]);
        decimals = parseInt(match[3]);
      }
    } else if (value) {
      const v3 = value;
      const check = (key, type, defaultValue) => {
        if (v3[key] == null) {
          return defaultValue;
        }
        assertArgument(typeof v3[key] === type, "invalid fixed format (" + key + " not " + type + ")", "format." + key, v3[key]);
        return v3[key];
      };
      signed2 = check("signed", "boolean", signed2);
      width = check("width", "number", width);
      decimals = check("decimals", "number", decimals);
    }
    assertArgument(width % 8 === 0, "invalid FixedNumber width (not byte aligned)", "format.width", width);
    assertArgument(decimals <= 80, "invalid FixedNumber decimals (too large)", "format.decimals", decimals);
    const name = (signed2 ? "" : "u") + "fixed" + String(width) + "x" + String(decimals);
    return { signed: signed2, width, decimals, name };
  }
  function toString(val, decimals) {
    let negative = "";
    if (val < BN_03) {
      negative = "-";
      val *= BN_N1;
    }
    let str = val.toString();
    if (decimals === 0) {
      return negative + str;
    }
    while (str.length <= decimals) {
      str = Zeros + str;
    }
    const index = str.length - decimals;
    str = str.substring(0, index) + "." + str.substring(index);
    while (str[0] === "0" && str[1] !== ".") {
      str = str.substring(1);
    }
    while (str[str.length - 1] === "0" && str[str.length - 2] !== ".") {
      str = str.substring(0, str.length - 1);
    }
    return negative + str;
  }
  var FixedNumber = class _FixedNumber {
    /**
     *  The specific fixed-point arithmetic field for this value.
     */
    format;
    #format;
    // The actual value (accounting for decimals)
    #val;
    // A base-10 value to multiple values by to maintain the magnitude
    #tens;
    /**
     *  This is a property so console.log shows a human-meaningful value.
     *
     *  @private
     */
    _value;
    // Use this when changing this file to get some typing info,
    // but then switch to any to mask the internal type
    //constructor(guard: any, value: bigint, format: _FixedFormat) {
    /**
     *  @private
     */
    constructor(guard, value, format) {
      assertPrivate(guard, _guard, "FixedNumber");
      this.#val = value;
      this.#format = format;
      const _value = toString(value, format.decimals);
      defineProperties(this, { format: format.name, _value });
      this.#tens = getTens(format.decimals);
    }
    /**
     *  If true, negative values are permitted, otherwise only
     *  positive values and zero are allowed.
     */
    get signed() {
      return this.#format.signed;
    }
    /**
     *  The number of bits available to store the value.
     */
    get width() {
      return this.#format.width;
    }
    /**
     *  The number of decimal places in the fixed-point arithment field.
     */
    get decimals() {
      return this.#format.decimals;
    }
    /**
     *  The value as an integer, based on the smallest unit the
     *  [[decimals]] allow.
     */
    get value() {
      return this.#val;
    }
    #checkFormat(other) {
      assertArgument(this.format === other.format, "incompatible format; use fixedNumber.toFormat", "other", other);
    }
    #checkValue(val, safeOp) {
      val = checkValue(val, this.#format, safeOp);
      return new _FixedNumber(_guard, val, this.#format);
    }
    #add(o4, safeOp) {
      this.#checkFormat(o4);
      return this.#checkValue(this.#val + o4.#val, safeOp);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% added
     *  to %%other%%, ignoring overflow.
     */
    addUnsafe(other) {
      return this.#add(other);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% added
     *  to %%other%%. A [[NumericFaultError]] is thrown if overflow
     *  occurs.
     */
    add(other) {
      return this.#add(other, "add");
    }
    #sub(o4, safeOp) {
      this.#checkFormat(o4);
      return this.#checkValue(this.#val - o4.#val, safeOp);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%other%% subtracted
     *  from %%this%%, ignoring overflow.
     */
    subUnsafe(other) {
      return this.#sub(other);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%other%% subtracted
     *  from %%this%%. A [[NumericFaultError]] is thrown if overflow
     *  occurs.
     */
    sub(other) {
      return this.#sub(other, "sub");
    }
    #mul(o4, safeOp) {
      this.#checkFormat(o4);
      return this.#checkValue(this.#val * o4.#val / this.#tens, safeOp);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% multiplied
     *  by %%other%%, ignoring overflow and underflow (precision loss).
     */
    mulUnsafe(other) {
      return this.#mul(other);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% multiplied
     *  by %%other%%. A [[NumericFaultError]] is thrown if overflow
     *  occurs.
     */
    mul(other) {
      return this.#mul(other, "mul");
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% multiplied
     *  by %%other%%. A [[NumericFaultError]] is thrown if overflow
     *  occurs or if underflow (precision loss) occurs.
     */
    mulSignal(other) {
      this.#checkFormat(other);
      const value = this.#val * other.#val;
      assert(value % this.#tens === BN_03, "precision lost during signalling mul", "NUMERIC_FAULT", {
        operation: "mulSignal",
        fault: "underflow",
        value: this
      });
      return this.#checkValue(value / this.#tens, "mulSignal");
    }
    #div(o4, safeOp) {
      assert(o4.#val !== BN_03, "division by zero", "NUMERIC_FAULT", {
        operation: "div",
        fault: "divide-by-zero",
        value: this
      });
      this.#checkFormat(o4);
      return this.#checkValue(this.#val * this.#tens / o4.#val, safeOp);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% divided
     *  by %%other%%, ignoring underflow (precision loss). A
     *  [[NumericFaultError]] is thrown if overflow occurs.
     */
    divUnsafe(other) {
      return this.#div(other);
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% divided
     *  by %%other%%, ignoring underflow (precision loss). A
     *  [[NumericFaultError]] is thrown if overflow occurs.
     */
    div(other) {
      return this.#div(other, "div");
    }
    /**
     *  Returns a new [[FixedNumber]] with the result of %%this%% divided
     *  by %%other%%. A [[NumericFaultError]] is thrown if underflow
     *  (precision loss) occurs.
     */
    divSignal(other) {
      assert(other.#val !== BN_03, "division by zero", "NUMERIC_FAULT", {
        operation: "div",
        fault: "divide-by-zero",
        value: this
      });
      this.#checkFormat(other);
      const value = this.#val * this.#tens;
      assert(value % other.#val === BN_03, "precision lost during signalling div", "NUMERIC_FAULT", {
        operation: "divSignal",
        fault: "underflow",
        value: this
      });
      return this.#checkValue(value / other.#val, "divSignal");
    }
    /**
     *  Returns a comparison result between %%this%% and %%other%%.
     *
     *  This is suitable for use in sorting, where ``-1`` implies %%this%%
     *  is smaller, ``1`` implies %%this%% is larger and ``0`` implies
     *  both are equal.
     */
    cmp(other) {
      let a4 = this.value, b5 = other.value;
      const delta = this.decimals - other.decimals;
      if (delta > 0) {
        b5 *= getTens(delta);
      } else if (delta < 0) {
        a4 *= getTens(-delta);
      }
      if (a4 < b5) {
        return -1;
      }
      if (a4 > b5) {
        return 1;
      }
      return 0;
    }
    /**
     *  Returns true if %%other%% is equal to %%this%%.
     */
    eq(other) {
      return this.cmp(other) === 0;
    }
    /**
     *  Returns true if %%other%% is less than to %%this%%.
     */
    lt(other) {
      return this.cmp(other) < 0;
    }
    /**
     *  Returns true if %%other%% is less than or equal to %%this%%.
     */
    lte(other) {
      return this.cmp(other) <= 0;
    }
    /**
     *  Returns true if %%other%% is greater than to %%this%%.
     */
    gt(other) {
      return this.cmp(other) > 0;
    }
    /**
     *  Returns true if %%other%% is greater than or equal to %%this%%.
     */
    gte(other) {
      return this.cmp(other) >= 0;
    }
    /**
     *  Returns a new [[FixedNumber]] which is the largest **integer**
     *  that is less than or equal to %%this%%.
     *
     *  The decimal component of the result will always be ``0``.
     */
    floor() {
      let val = this.#val;
      if (this.#val < BN_03) {
        val -= this.#tens - BN_12;
      }
      val = this.#val / this.#tens * this.#tens;
      return this.#checkValue(val, "floor");
    }
    /**
     *  Returns a new [[FixedNumber]] which is the smallest **integer**
     *  that is greater than or equal to %%this%%.
     *
     *  The decimal component of the result will always be ``0``.
     */
    ceiling() {
      let val = this.#val;
      if (this.#val > BN_03) {
        val += this.#tens - BN_12;
      }
      val = this.#val / this.#tens * this.#tens;
      return this.#checkValue(val, "ceiling");
    }
    /**
     *  Returns a new [[FixedNumber]] with the decimal component
     *  rounded up on ties at %%decimals%% places.
     */
    round(decimals) {
      if (decimals == null) {
        decimals = 0;
      }
      if (decimals >= this.decimals) {
        return this;
      }
      const delta = this.decimals - decimals;
      const bump = BN_5 * getTens(delta - 1);
      let value = this.value + bump;
      const tens = getTens(delta);
      value = value / tens * tens;
      checkValue(value, this.#format, "round");
      return new _FixedNumber(_guard, value, this.#format);
    }
    /**
     *  Returns true if %%this%% is equal to ``0``.
     */
    isZero() {
      return this.#val === BN_03;
    }
    /**
     *  Returns true if %%this%% is less than ``0``.
     */
    isNegative() {
      return this.#val < BN_03;
    }
    /**
     *  Returns the string representation of %%this%%.
     */
    toString() {
      return this._value;
    }
    /**
     *  Returns a float approximation.
     *
     *  Due to IEEE 754 precission (or lack thereof), this function
     *  can only return an approximation and most values will contain
     *  rounding errors.
     */
    toUnsafeFloat() {
      return parseFloat(this.toString());
    }
    /**
     *  Return a new [[FixedNumber]] with the same value but has had
     *  its field set to %%format%%.
     *
     *  This will throw if the value cannot fit into %%format%%.
     */
    toFormat(format) {
      return _FixedNumber.fromString(this.toString(), format);
    }
    /**
     *  Creates a new [[FixedNumber]] for %%value%% divided by
     *  %%decimal%% places with %%format%%.
     *
     *  This will throw a [[NumericFaultError]] if %%value%% (once adjusted
     *  for %%decimals%%) cannot fit in %%format%%, either due to overflow
     *  or underflow (precision loss).
     */
    static fromValue(_value, _decimals, _format) {
      const decimals = _decimals == null ? 0 : getNumber(_decimals);
      const format = getFormat(_format);
      let value = getBigInt(_value, "value");
      const delta = decimals - format.decimals;
      if (delta > 0) {
        const tens = getTens(delta);
        assert(value % tens === BN_03, "value loses precision for format", "NUMERIC_FAULT", {
          operation: "fromValue",
          fault: "underflow",
          value: _value
        });
        value /= tens;
      } else if (delta < 0) {
        value *= getTens(-delta);
      }
      checkValue(value, format, "fromValue");
      return new _FixedNumber(_guard, value, format);
    }
    /**
     *  Creates a new [[FixedNumber]] for %%value%% with %%format%%.
     *
     *  This will throw a [[NumericFaultError]] if %%value%% cannot fit
     *  in %%format%%, either due to overflow or underflow (precision loss).
     */
    static fromString(_value, _format) {
      const match = _value.match(/^(-?)([0-9]*)\.?([0-9]*)$/);
      assertArgument(match && match[2].length + match[3].length > 0, "invalid FixedNumber string value", "value", _value);
      const format = getFormat(_format);
      let whole = match[2] || "0", decimal = match[3] || "";
      while (decimal.length < format.decimals) {
        decimal += Zeros;
      }
      assert(decimal.substring(format.decimals).match(/^0*$/), "too many decimals for format", "NUMERIC_FAULT", {
        operation: "fromString",
        fault: "underflow",
        value: _value
      });
      decimal = decimal.substring(0, format.decimals);
      const value = BigInt(match[1] + whole + decimal);
      checkValue(value, format, "fromString");
      return new _FixedNumber(_guard, value, format);
    }
    /**
     *  Creates a new [[FixedNumber]] with the big-endian representation
     *  %%value%% with %%format%%.
     *
     *  This will throw a [[NumericFaultError]] if %%value%% cannot fit
     *  in %%format%% due to overflow.
     */
    static fromBytes(_value, _format) {
      let value = toBigInt(getBytes(_value, "value"));
      const format = getFormat(_format);
      if (format.signed) {
        value = fromTwos(value, format.width);
      }
      checkValue(value, format, "fromBytes");
      return new _FixedNumber(_guard, value, format);
    }
  };

  // node_modules/ethers/lib.esm/utils/rlp-decode.js
  function hexlifyByte(value) {
    let result = value.toString(16);
    while (result.length < 2) {
      result = "0" + result;
    }
    return "0x" + result;
  }
  function unarrayifyInteger(data, offset, length) {
    let result = 0;
    for (let i4 = 0; i4 < length; i4++) {
      result = result * 256 + data[offset + i4];
    }
    return result;
  }
  function _decodeChildren(data, offset, childOffset, length) {
    const result = [];
    while (childOffset < offset + 1 + length) {
      const decoded = _decode(data, childOffset);
      result.push(decoded.result);
      childOffset += decoded.consumed;
      assert(childOffset <= offset + 1 + length, "child data too short", "BUFFER_OVERRUN", {
        buffer: data,
        length,
        offset
      });
    }
    return { consumed: 1 + length, result };
  }
  function _decode(data, offset) {
    assert(data.length !== 0, "data too short", "BUFFER_OVERRUN", {
      buffer: data,
      length: 0,
      offset: 1
    });
    const checkOffset = (offset2) => {
      assert(offset2 <= data.length, "data short segment too short", "BUFFER_OVERRUN", {
        buffer: data,
        length: data.length,
        offset: offset2
      });
    };
    if (data[offset] >= 248) {
      const lengthLength = data[offset] - 247;
      checkOffset(offset + 1 + lengthLength);
      const length = unarrayifyInteger(data, offset + 1, lengthLength);
      checkOffset(offset + 1 + lengthLength + length);
      return _decodeChildren(data, offset, offset + 1 + lengthLength, lengthLength + length);
    } else if (data[offset] >= 192) {
      const length = data[offset] - 192;
      checkOffset(offset + 1 + length);
      return _decodeChildren(data, offset, offset + 1, length);
    } else if (data[offset] >= 184) {
      const lengthLength = data[offset] - 183;
      checkOffset(offset + 1 + lengthLength);
      const length = unarrayifyInteger(data, offset + 1, lengthLength);
      checkOffset(offset + 1 + lengthLength + length);
      const result = hexlify(data.slice(offset + 1 + lengthLength, offset + 1 + lengthLength + length));
      return { consumed: 1 + lengthLength + length, result };
    } else if (data[offset] >= 128) {
      const length = data[offset] - 128;
      checkOffset(offset + 1 + length);
      const result = hexlify(data.slice(offset + 1, offset + 1 + length));
      return { consumed: 1 + length, result };
    }
    return { consumed: 1, result: hexlifyByte(data[offset]) };
  }
  function decodeRlp(_data) {
    const data = getBytes(_data, "data");
    const decoded = _decode(data, 0);
    assertArgument(decoded.consumed === data.length, "unexpected junk after rlp payload", "data", _data);
    return decoded.result;
  }

  // node_modules/ethers/lib.esm/utils/rlp-encode.js
  function arrayifyInteger(value) {
    const result = [];
    while (value) {
      result.unshift(value & 255);
      value >>= 8;
    }
    return result;
  }
  function _encode(object2) {
    if (Array.isArray(object2)) {
      let payload = [];
      object2.forEach(function(child) {
        payload = payload.concat(_encode(child));
      });
      if (payload.length <= 55) {
        payload.unshift(192 + payload.length);
        return payload;
      }
      const length2 = arrayifyInteger(payload.length);
      length2.unshift(247 + length2.length);
      return length2.concat(payload);
    }
    const data = Array.prototype.slice.call(getBytes(object2, "object"));
    if (data.length === 1 && data[0] <= 127) {
      return data;
    } else if (data.length <= 55) {
      data.unshift(128 + data.length);
      return data;
    }
    const length = arrayifyInteger(data.length);
    length.unshift(183 + length.length);
    return length.concat(data);
  }
  var nibbles = "0123456789abcdef";
  function encodeRlp(object2) {
    let result = "0x";
    for (const v3 of _encode(object2)) {
      result += nibbles[v3 >> 4];
      result += nibbles[v3 & 15];
    }
    return result;
  }

  // node_modules/ethers/lib.esm/utils/units.js
  var names = [
    "wei",
    "kwei",
    "mwei",
    "gwei",
    "szabo",
    "finney",
    "ether"
  ];
  function formatUnits(value, unit) {
    let decimals = 18;
    if (typeof unit === "string") {
      const index = names.indexOf(unit);
      assertArgument(index >= 0, "invalid unit", "unit", unit);
      decimals = 3 * index;
    } else if (unit != null) {
      decimals = getNumber(unit, "unit");
    }
    return FixedNumber.fromValue(value, decimals, { decimals, width: 512 }).toString();
  }
  function parseUnits(value, unit) {
    assertArgument(typeof value === "string", "value must be a string", "value", value);
    let decimals = 18;
    if (typeof unit === "string") {
      const index = names.indexOf(unit);
      assertArgument(index >= 0, "invalid unit", "unit", unit);
      decimals = 3 * index;
    } else if (unit != null) {
      decimals = getNumber(unit, "unit");
    }
    return FixedNumber.fromString(value, { decimals, width: 512 }).value;
  }
  function formatEther(wei) {
    return formatUnits(wei, 18);
  }

  // node_modules/ethers/lib.esm/abi/coders/abstract-coder.js
  var WordSize = 32;
  var Padding = new Uint8Array(WordSize);
  var passProperties = ["then"];
  var _guard2 = {};
  function throwError(name, error) {
    const wrapped = new Error(`deferred error during ABI decoding triggered accessing ${name}`);
    wrapped.error = error;
    throw wrapped;
  }
  var Result = class _Result extends Array {
    #names;
    /**
     *  @private
     */
    constructor(...args) {
      const guard = args[0];
      let items = args[1];
      let names2 = (args[2] || []).slice();
      let wrap = true;
      if (guard !== _guard2) {
        items = args;
        names2 = [];
        wrap = false;
      }
      super(items.length);
      items.forEach((item, index) => {
        this[index] = item;
      });
      const nameCounts = names2.reduce((accum, name) => {
        if (typeof name === "string") {
          accum.set(name, (accum.get(name) || 0) + 1);
        }
        return accum;
      }, /* @__PURE__ */ new Map());
      this.#names = Object.freeze(items.map((item, index) => {
        const name = names2[index];
        if (name != null && nameCounts.get(name) === 1) {
          return name;
        }
        return null;
      }));
      if (!wrap) {
        return;
      }
      Object.freeze(this);
      return new Proxy(this, {
        get: (target, prop, receiver) => {
          if (typeof prop === "string") {
            if (prop.match(/^[0-9]+$/)) {
              const index = getNumber(prop, "%index");
              if (index < 0 || index >= this.length) {
                throw new RangeError("out of result range");
              }
              const item = target[index];
              if (item instanceof Error) {
                throwError(`index ${index}`, item);
              }
              return item;
            }
            if (passProperties.indexOf(prop) >= 0) {
              return Reflect.get(target, prop, receiver);
            }
            const value = target[prop];
            if (value instanceof Function) {
              return function(...args2) {
                return value.apply(this === receiver ? target : this, args2);
              };
            } else if (!(prop in target)) {
              return target.getValue.apply(this === receiver ? target : this, [prop]);
            }
          }
          return Reflect.get(target, prop, receiver);
        }
      });
    }
    /**
     *  Returns the Result as a normal Array.
     *
     *  This will throw if there are any outstanding deferred
     *  errors.
     */
    toArray() {
      const result = [];
      this.forEach((item, index) => {
        if (item instanceof Error) {
          throwError(`index ${index}`, item);
        }
        result.push(item);
      });
      return result;
    }
    /**
     *  Returns the Result as an Object with each name-value pair.
     *
     *  This will throw if any value is unnamed, or if there are
     *  any outstanding deferred errors.
     */
    toObject() {
      return this.#names.reduce((accum, name, index) => {
        assert(name != null, "value at index ${ index } unnamed", "UNSUPPORTED_OPERATION", {
          operation: "toObject()"
        });
        if (!(name in accum)) {
          accum[name] = this.getValue(name);
        }
        return accum;
      }, {});
    }
    /**
     *  @_ignore
     */
    slice(start, end) {
      if (start == null) {
        start = 0;
      }
      if (start < 0) {
        start += this.length;
        if (start < 0) {
          start = 0;
        }
      }
      if (end == null) {
        end = this.length;
      }
      if (end < 0) {
        end += this.length;
        if (end < 0) {
          end = 0;
        }
      }
      if (end > this.length) {
        end = this.length;
      }
      const result = [], names2 = [];
      for (let i4 = start; i4 < end; i4++) {
        result.push(this[i4]);
        names2.push(this.#names[i4]);
      }
      return new _Result(_guard2, result, names2);
    }
    /**
     *  @_ignore
     */
    filter(callback, thisArg) {
      const result = [], names2 = [];
      for (let i4 = 0; i4 < this.length; i4++) {
        const item = this[i4];
        if (item instanceof Error) {
          throwError(`index ${i4}`, item);
        }
        if (callback.call(thisArg, item, i4, this)) {
          result.push(item);
          names2.push(this.#names[i4]);
        }
      }
      return new _Result(_guard2, result, names2);
    }
    /**
     *  @_ignore
     */
    map(callback, thisArg) {
      const result = [];
      for (let i4 = 0; i4 < this.length; i4++) {
        const item = this[i4];
        if (item instanceof Error) {
          throwError(`index ${i4}`, item);
        }
        result.push(callback.call(thisArg, item, i4, this));
      }
      return result;
    }
    /**
     *  Returns the value for %%name%%.
     *
     *  Since it is possible to have a key whose name conflicts with
     *  a method on a [[Result]] or its superclass Array, or any
     *  JavaScript keyword, this ensures all named values are still
     *  accessible by name.
     */
    getValue(name) {
      const index = this.#names.indexOf(name);
      if (index === -1) {
        return void 0;
      }
      const value = this[index];
      if (value instanceof Error) {
        throwError(`property ${JSON.stringify(name)}`, value.error);
      }
      return value;
    }
    /**
     *  Creates a new [[Result]] for %%items%% with each entry
     *  also accessible by its corresponding name in %%keys%%.
     */
    static fromItems(items, keys) {
      return new _Result(_guard2, items, keys);
    }
  };
  function getValue(value) {
    let bytes2 = toBeArray(value);
    assert(bytes2.length <= WordSize, "value out-of-bounds", "BUFFER_OVERRUN", { buffer: bytes2, length: WordSize, offset: bytes2.length });
    if (bytes2.length !== WordSize) {
      bytes2 = getBytesCopy(concat([Padding.slice(bytes2.length % WordSize), bytes2]));
    }
    return bytes2;
  }
  var Coder = class {
    // The coder name:
    //   - address, uint256, tuple, array, etc.
    name;
    // The fully expanded type, including composite types:
    //   - address, uint256, tuple(address,bytes), uint256[3][4][],  etc.
    type;
    // The localName bound in the signature, in this example it is "baz":
    //   - tuple(address foo, uint bar) baz
    localName;
    // Whether this type is dynamic:
    //  - Dynamic: bytes, string, address[], tuple(boolean[]), etc.
    //  - Not Dynamic: address, uint256, boolean[3], tuple(address, uint8)
    dynamic;
    constructor(name, type, localName, dynamic) {
      defineProperties(this, { name, type, localName, dynamic }, {
        name: "string",
        type: "string",
        localName: "string",
        dynamic: "boolean"
      });
    }
    _throwError(message, value) {
      assertArgument(false, message, this.localName, value);
    }
  };
  var Writer = class {
    // An array of WordSize lengthed objects to concatenation
    #data;
    #dataLength;
    constructor() {
      this.#data = [];
      this.#dataLength = 0;
    }
    get data() {
      return concat(this.#data);
    }
    get length() {
      return this.#dataLength;
    }
    #writeData(data) {
      this.#data.push(data);
      this.#dataLength += data.length;
      return data.length;
    }
    appendWriter(writer) {
      return this.#writeData(getBytesCopy(writer.data));
    }
    // Arrayish item; pad on the right to *nearest* WordSize
    writeBytes(value) {
      let bytes2 = getBytesCopy(value);
      const paddingOffset = bytes2.length % WordSize;
      if (paddingOffset) {
        bytes2 = getBytesCopy(concat([bytes2, Padding.slice(paddingOffset)]));
      }
      return this.#writeData(bytes2);
    }
    // Numeric item; pad on the left *to* WordSize
    writeValue(value) {
      return this.#writeData(getValue(value));
    }
    // Inserts a numeric place-holder, returning a callback that can
    // be used to asjust the value later
    writeUpdatableValue() {
      const offset = this.#data.length;
      this.#data.push(Padding);
      this.#dataLength += WordSize;
      return (value) => {
        this.#data[offset] = getValue(value);
      };
    }
  };
  var Reader = class _Reader {
    // Allows incomplete unpadded data to be read; otherwise an error
    // is raised if attempting to overrun the buffer. This is required
    // to deal with an old Solidity bug, in which event data for
    // external (not public thoguh) was tightly packed.
    allowLoose;
    #data;
    #offset;
    constructor(data, allowLoose) {
      defineProperties(this, { allowLoose: !!allowLoose });
      this.#data = getBytesCopy(data);
      this.#offset = 0;
    }
    get data() {
      return hexlify(this.#data);
    }
    get dataLength() {
      return this.#data.length;
    }
    get consumed() {
      return this.#offset;
    }
    get bytes() {
      return new Uint8Array(this.#data);
    }
    #peekBytes(offset, length, loose) {
      let alignedLength = Math.ceil(length / WordSize) * WordSize;
      if (this.#offset + alignedLength > this.#data.length) {
        if (this.allowLoose && loose && this.#offset + length <= this.#data.length) {
          alignedLength = length;
        } else {
          assert(false, "data out-of-bounds", "BUFFER_OVERRUN", {
            buffer: getBytesCopy(this.#data),
            length: this.#data.length,
            offset: this.#offset + alignedLength
          });
        }
      }
      return this.#data.slice(this.#offset, this.#offset + alignedLength);
    }
    // Create a sub-reader with the same underlying data, but offset
    subReader(offset) {
      return new _Reader(this.#data.slice(this.#offset + offset), this.allowLoose);
    }
    // Read bytes
    readBytes(length, loose) {
      let bytes2 = this.#peekBytes(0, length, !!loose);
      this.#offset += bytes2.length;
      return bytes2.slice(0, length);
    }
    // Read a numeric values
    readValue() {
      return toBigInt(this.readBytes(WordSize));
    }
    readIndex() {
      return toNumber(this.readBytes(WordSize));
    }
  };

  // node_modules/@noble/hashes/esm/_assert.js
  function number(n4) {
    if (!Number.isSafeInteger(n4) || n4 < 0)
      throw new Error(`Wrong positive integer: ${n4}`);
  }
  function bool(b5) {
    if (typeof b5 !== "boolean")
      throw new Error(`Expected boolean, not ${b5}`);
  }
  function bytes(b5, ...lengths) {
    if (!(b5 instanceof Uint8Array))
      throw new TypeError("Expected Uint8Array");
    if (lengths.length > 0 && !lengths.includes(b5.length))
      throw new TypeError(`Expected Uint8Array of length ${lengths}, not of length=${b5.length}`);
  }
  function hash(hash2) {
    if (typeof hash2 !== "function" || typeof hash2.create !== "function")
      throw new Error("Hash should be wrapped by utils.wrapConstructor");
    number(hash2.outputLen);
    number(hash2.blockLen);
  }
  function exists(instance, checkFinished = true) {
    if (instance.destroyed)
      throw new Error("Hash instance has been destroyed");
    if (checkFinished && instance.finished)
      throw new Error("Hash#digest() has already been called");
  }
  function output(out, instance) {
    bytes(out);
    const min = instance.outputLen;
    if (out.length < min) {
      throw new Error(`digestInto() expects output buffer of length at least ${min}`);
    }
  }
  var assert2 = {
    number,
    bool,
    bytes,
    hash,
    exists,
    output
  };
  var assert_default = assert2;

  // node_modules/@noble/hashes/esm/cryptoBrowser.js
  var crypto = {
    node: void 0,
    web: typeof self === "object" && "crypto" in self ? self.crypto : void 0
  };

  // node_modules/@noble/hashes/esm/utils.js
  var u32 = (arr) => new Uint32Array(arr.buffer, arr.byteOffset, Math.floor(arr.byteLength / 4));
  var createView = (arr) => new DataView(arr.buffer, arr.byteOffset, arr.byteLength);
  var rotr = (word, shift) => word << 32 - shift | word >>> shift;
  var isLE = new Uint8Array(new Uint32Array([287454020]).buffer)[0] === 68;
  if (!isLE)
    throw new Error("Non little-endian hardware is not supported");
  var hexes = Array.from({ length: 256 }, (v3, i4) => i4.toString(16).padStart(2, "0"));
  function utf8ToBytes(str) {
    if (typeof str !== "string") {
      throw new TypeError(`utf8ToBytes expected string, got ${typeof str}`);
    }
    return new TextEncoder().encode(str);
  }
  function toBytes(data) {
    if (typeof data === "string")
      data = utf8ToBytes(data);
    if (!(data instanceof Uint8Array))
      throw new TypeError(`Expected input type is Uint8Array (got ${typeof data})`);
    return data;
  }
  var Hash = class {
    // Safe version that clones internal state
    clone() {
      return this._cloneInto();
    }
  };
  function wrapConstructor(hashConstructor) {
    const hashC = (message) => hashConstructor().update(toBytes(message)).digest();
    const tmp = hashConstructor();
    hashC.outputLen = tmp.outputLen;
    hashC.blockLen = tmp.blockLen;
    hashC.create = () => hashConstructor();
    return hashC;
  }
  function wrapConstructorWithOpts(hashCons) {
    const hashC = (msg, opts) => hashCons(opts).update(toBytes(msg)).digest();
    const tmp = hashCons({});
    hashC.outputLen = tmp.outputLen;
    hashC.blockLen = tmp.blockLen;
    hashC.create = (opts) => hashCons(opts);
    return hashC;
  }

  // node_modules/@noble/hashes/esm/hmac.js
  var HMAC = class extends Hash {
    constructor(hash2, _key) {
      super();
      this.finished = false;
      this.destroyed = false;
      assert_default.hash(hash2);
      const key = toBytes(_key);
      this.iHash = hash2.create();
      if (!(this.iHash instanceof Hash))
        throw new TypeError("Expected instance of class which extends utils.Hash");
      const blockLen = this.blockLen = this.iHash.blockLen;
      this.outputLen = this.iHash.outputLen;
      const pad = new Uint8Array(blockLen);
      pad.set(key.length > this.iHash.blockLen ? hash2.create().update(key).digest() : key);
      for (let i4 = 0; i4 < pad.length; i4++)
        pad[i4] ^= 54;
      this.iHash.update(pad);
      this.oHash = hash2.create();
      for (let i4 = 0; i4 < pad.length; i4++)
        pad[i4] ^= 54 ^ 92;
      this.oHash.update(pad);
      pad.fill(0);
    }
    update(buf) {
      assert_default.exists(this);
      this.iHash.update(buf);
      return this;
    }
    digestInto(out) {
      assert_default.exists(this);
      assert_default.bytes(out, this.outputLen);
      this.finished = true;
      this.iHash.digestInto(out);
      this.oHash.update(out);
      this.oHash.digestInto(out);
      this.destroy();
    }
    digest() {
      const out = new Uint8Array(this.oHash.outputLen);
      this.digestInto(out);
      return out;
    }
    _cloneInto(to) {
      to || (to = Object.create(Object.getPrototypeOf(this), {}));
      const { oHash, iHash, finished, destroyed, blockLen, outputLen } = this;
      to = to;
      to.finished = finished;
      to.destroyed = destroyed;
      to.blockLen = blockLen;
      to.outputLen = outputLen;
      to.oHash = oHash._cloneInto(to.oHash);
      to.iHash = iHash._cloneInto(to.iHash);
      return to;
    }
    destroy() {
      this.destroyed = true;
      this.oHash.destroy();
      this.iHash.destroy();
    }
  };
  var hmac = (hash2, key, message) => new HMAC(hash2, key).update(message).digest();
  hmac.create = (hash2, key) => new HMAC(hash2, key);

  // node_modules/@noble/hashes/esm/_sha2.js
  function setBigUint64(view, byteOffset, value, isLE2) {
    if (typeof view.setBigUint64 === "function")
      return view.setBigUint64(byteOffset, value, isLE2);
    const _32n2 = BigInt(32);
    const _u32_max = BigInt(4294967295);
    const wh = Number(value >> _32n2 & _u32_max);
    const wl = Number(value & _u32_max);
    const h4 = isLE2 ? 4 : 0;
    const l4 = isLE2 ? 0 : 4;
    view.setUint32(byteOffset + h4, wh, isLE2);
    view.setUint32(byteOffset + l4, wl, isLE2);
  }
  var SHA2 = class extends Hash {
    constructor(blockLen, outputLen, padOffset, isLE2) {
      super();
      this.blockLen = blockLen;
      this.outputLen = outputLen;
      this.padOffset = padOffset;
      this.isLE = isLE2;
      this.finished = false;
      this.length = 0;
      this.pos = 0;
      this.destroyed = false;
      this.buffer = new Uint8Array(blockLen);
      this.view = createView(this.buffer);
    }
    update(data) {
      assert_default.exists(this);
      const { view, buffer, blockLen } = this;
      data = toBytes(data);
      const len = data.length;
      for (let pos = 0; pos < len; ) {
        const take = Math.min(blockLen - this.pos, len - pos);
        if (take === blockLen) {
          const dataView = createView(data);
          for (; blockLen <= len - pos; pos += blockLen)
            this.process(dataView, pos);
          continue;
        }
        buffer.set(data.subarray(pos, pos + take), this.pos);
        this.pos += take;
        pos += take;
        if (this.pos === blockLen) {
          this.process(view, 0);
          this.pos = 0;
        }
      }
      this.length += data.length;
      this.roundClean();
      return this;
    }
    digestInto(out) {
      assert_default.exists(this);
      assert_default.output(out, this);
      this.finished = true;
      const { buffer, view, blockLen, isLE: isLE2 } = this;
      let { pos } = this;
      buffer[pos++] = 128;
      this.buffer.subarray(pos).fill(0);
      if (this.padOffset > blockLen - pos) {
        this.process(view, 0);
        pos = 0;
      }
      for (let i4 = pos; i4 < blockLen; i4++)
        buffer[i4] = 0;
      setBigUint64(view, blockLen - 8, BigInt(this.length * 8), isLE2);
      this.process(view, 0);
      const oview = createView(out);
      this.get().forEach((v3, i4) => oview.setUint32(4 * i4, v3, isLE2));
    }
    digest() {
      const { buffer, outputLen } = this;
      this.digestInto(buffer);
      const res = buffer.slice(0, outputLen);
      this.destroy();
      return res;
    }
    _cloneInto(to) {
      to || (to = new this.constructor());
      to.set(...this.get());
      const { blockLen, buffer, length, finished, destroyed, pos } = this;
      to.length = length;
      to.pos = pos;
      to.finished = finished;
      to.destroyed = destroyed;
      if (length % blockLen)
        to.buffer.set(buffer);
      return to;
    }
  };

  // node_modules/@noble/hashes/esm/sha256.js
  var Chi = (a4, b5, c4) => a4 & b5 ^ ~a4 & c4;
  var Maj = (a4, b5, c4) => a4 & b5 ^ a4 & c4 ^ b5 & c4;
  var SHA256_K = new Uint32Array([
    1116352408,
    1899447441,
    3049323471,
    3921009573,
    961987163,
    1508970993,
    2453635748,
    2870763221,
    3624381080,
    310598401,
    607225278,
    1426881987,
    1925078388,
    2162078206,
    2614888103,
    3248222580,
    3835390401,
    4022224774,
    264347078,
    604807628,
    770255983,
    1249150122,
    1555081692,
    1996064986,
    2554220882,
    2821834349,
    2952996808,
    3210313671,
    3336571891,
    3584528711,
    113926993,
    338241895,
    666307205,
    773529912,
    1294757372,
    1396182291,
    1695183700,
    1986661051,
    2177026350,
    2456956037,
    2730485921,
    2820302411,
    3259730800,
    3345764771,
    3516065817,
    3600352804,
    4094571909,
    275423344,
    430227734,
    506948616,
    659060556,
    883997877,
    958139571,
    1322822218,
    1537002063,
    1747873779,
    1955562222,
    2024104815,
    2227730452,
    2361852424,
    2428436474,
    2756734187,
    3204031479,
    3329325298
  ]);
  var IV = new Uint32Array([
    1779033703,
    3144134277,
    1013904242,
    2773480762,
    1359893119,
    2600822924,
    528734635,
    1541459225
  ]);
  var SHA256_W = new Uint32Array(64);
  var SHA256 = class extends SHA2 {
    constructor() {
      super(64, 32, 8, false);
      this.A = IV[0] | 0;
      this.B = IV[1] | 0;
      this.C = IV[2] | 0;
      this.D = IV[3] | 0;
      this.E = IV[4] | 0;
      this.F = IV[5] | 0;
      this.G = IV[6] | 0;
      this.H = IV[7] | 0;
    }
    get() {
      const { A: A2, B: B3, C, D: D2, E, F: F2, G, H: H2 } = this;
      return [A2, B3, C, D2, E, F2, G, H2];
    }
    // prettier-ignore
    set(A2, B3, C, D2, E, F2, G, H2) {
      this.A = A2 | 0;
      this.B = B3 | 0;
      this.C = C | 0;
      this.D = D2 | 0;
      this.E = E | 0;
      this.F = F2 | 0;
      this.G = G | 0;
      this.H = H2 | 0;
    }
    process(view, offset) {
      for (let i4 = 0; i4 < 16; i4++, offset += 4)
        SHA256_W[i4] = view.getUint32(offset, false);
      for (let i4 = 16; i4 < 64; i4++) {
        const W15 = SHA256_W[i4 - 15];
        const W2 = SHA256_W[i4 - 2];
        const s0 = rotr(W15, 7) ^ rotr(W15, 18) ^ W15 >>> 3;
        const s1 = rotr(W2, 17) ^ rotr(W2, 19) ^ W2 >>> 10;
        SHA256_W[i4] = s1 + SHA256_W[i4 - 7] + s0 + SHA256_W[i4 - 16] | 0;
      }
      let { A: A2, B: B3, C, D: D2, E, F: F2, G, H: H2 } = this;
      for (let i4 = 0; i4 < 64; i4++) {
        const sigma1 = rotr(E, 6) ^ rotr(E, 11) ^ rotr(E, 25);
        const T12 = H2 + sigma1 + Chi(E, F2, G) + SHA256_K[i4] + SHA256_W[i4] | 0;
        const sigma0 = rotr(A2, 2) ^ rotr(A2, 13) ^ rotr(A2, 22);
        const T22 = sigma0 + Maj(A2, B3, C) | 0;
        H2 = G;
        G = F2;
        F2 = E;
        E = D2 + T12 | 0;
        D2 = C;
        C = B3;
        B3 = A2;
        A2 = T12 + T22 | 0;
      }
      A2 = A2 + this.A | 0;
      B3 = B3 + this.B | 0;
      C = C + this.C | 0;
      D2 = D2 + this.D | 0;
      E = E + this.E | 0;
      F2 = F2 + this.F | 0;
      G = G + this.G | 0;
      H2 = H2 + this.H | 0;
      this.set(A2, B3, C, D2, E, F2, G, H2);
    }
    roundClean() {
      SHA256_W.fill(0);
    }
    destroy() {
      this.set(0, 0, 0, 0, 0, 0, 0, 0);
      this.buffer.fill(0);
    }
  };
  var sha256 = wrapConstructor(() => new SHA256());

  // node_modules/@noble/hashes/esm/_u64.js
  var U32_MASK64 = BigInt(2 ** 32 - 1);
  var _32n = BigInt(32);
  function fromBig(n4, le = false) {
    if (le)
      return { h: Number(n4 & U32_MASK64), l: Number(n4 >> _32n & U32_MASK64) };
    return { h: Number(n4 >> _32n & U32_MASK64) | 0, l: Number(n4 & U32_MASK64) | 0 };
  }
  function split(lst, le = false) {
    let Ah = new Uint32Array(lst.length);
    let Al = new Uint32Array(lst.length);
    for (let i4 = 0; i4 < lst.length; i4++) {
      const { h: h4, l: l4 } = fromBig(lst[i4], le);
      [Ah[i4], Al[i4]] = [h4, l4];
    }
    return [Ah, Al];
  }
  var toBig = (h4, l4) => BigInt(h4 >>> 0) << _32n | BigInt(l4 >>> 0);
  var shrSH = (h4, l4, s4) => h4 >>> s4;
  var shrSL = (h4, l4, s4) => h4 << 32 - s4 | l4 >>> s4;
  var rotrSH = (h4, l4, s4) => h4 >>> s4 | l4 << 32 - s4;
  var rotrSL = (h4, l4, s4) => h4 << 32 - s4 | l4 >>> s4;
  var rotrBH = (h4, l4, s4) => h4 << 64 - s4 | l4 >>> s4 - 32;
  var rotrBL = (h4, l4, s4) => h4 >>> s4 - 32 | l4 << 64 - s4;
  var rotr32H = (h4, l4) => l4;
  var rotr32L = (h4, l4) => h4;
  var rotlSH = (h4, l4, s4) => h4 << s4 | l4 >>> 32 - s4;
  var rotlSL = (h4, l4, s4) => l4 << s4 | h4 >>> 32 - s4;
  var rotlBH = (h4, l4, s4) => l4 << s4 - 32 | h4 >>> 64 - s4;
  var rotlBL = (h4, l4, s4) => h4 << s4 - 32 | l4 >>> 64 - s4;
  function add(Ah, Al, Bh, Bl) {
    const l4 = (Al >>> 0) + (Bl >>> 0);
    return { h: Ah + Bh + (l4 / 2 ** 32 | 0) | 0, l: l4 | 0 };
  }
  var add3L = (Al, Bl, Cl) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0);
  var add3H = (low, Ah, Bh, Ch) => Ah + Bh + Ch + (low / 2 ** 32 | 0) | 0;
  var add4L = (Al, Bl, Cl, Dl) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0) + (Dl >>> 0);
  var add4H = (low, Ah, Bh, Ch, Dh) => Ah + Bh + Ch + Dh + (low / 2 ** 32 | 0) | 0;
  var add5L = (Al, Bl, Cl, Dl, El) => (Al >>> 0) + (Bl >>> 0) + (Cl >>> 0) + (Dl >>> 0) + (El >>> 0);
  var add5H = (low, Ah, Bh, Ch, Dh, Eh) => Ah + Bh + Ch + Dh + Eh + (low / 2 ** 32 | 0) | 0;
  var u64 = {
    fromBig,
    split,
    toBig,
    shrSH,
    shrSL,
    rotrSH,
    rotrSL,
    rotrBH,
    rotrBL,
    rotr32H,
    rotr32L,
    rotlSH,
    rotlSL,
    rotlBH,
    rotlBL,
    add,
    add3L,
    add3H,
    add4L,
    add4H,
    add5H,
    add5L
  };
  var u64_default = u64;

  // node_modules/@noble/hashes/esm/sha512.js
  var [SHA512_Kh, SHA512_Kl] = u64_default.split([
    "0x428a2f98d728ae22",
    "0x7137449123ef65cd",
    "0xb5c0fbcfec4d3b2f",
    "0xe9b5dba58189dbbc",
    "0x3956c25bf348b538",
    "0x59f111f1b605d019",
    "0x923f82a4af194f9b",
    "0xab1c5ed5da6d8118",
    "0xd807aa98a3030242",
    "0x12835b0145706fbe",
    "0x243185be4ee4b28c",
    "0x550c7dc3d5ffb4e2",
    "0x72be5d74f27b896f",
    "0x80deb1fe3b1696b1",
    "0x9bdc06a725c71235",
    "0xc19bf174cf692694",
    "0xe49b69c19ef14ad2",
    "0xefbe4786384f25e3",
    "0x0fc19dc68b8cd5b5",
    "0x240ca1cc77ac9c65",
    "0x2de92c6f592b0275",
    "0x4a7484aa6ea6e483",
    "0x5cb0a9dcbd41fbd4",
    "0x76f988da831153b5",
    "0x983e5152ee66dfab",
    "0xa831c66d2db43210",
    "0xb00327c898fb213f",
    "0xbf597fc7beef0ee4",
    "0xc6e00bf33da88fc2",
    "0xd5a79147930aa725",
    "0x06ca6351e003826f",
    "0x142929670a0e6e70",
    "0x27b70a8546d22ffc",
    "0x2e1b21385c26c926",
    "0x4d2c6dfc5ac42aed",
    "0x53380d139d95b3df",
    "0x650a73548baf63de",
    "0x766a0abb3c77b2a8",
    "0x81c2c92e47edaee6",
    "0x92722c851482353b",
    "0xa2bfe8a14cf10364",
    "0xa81a664bbc423001",
    "0xc24b8b70d0f89791",
    "0xc76c51a30654be30",
    "0xd192e819d6ef5218",
    "0xd69906245565a910",
    "0xf40e35855771202a",
    "0x106aa07032bbd1b8",
    "0x19a4c116b8d2d0c8",
    "0x1e376c085141ab53",
    "0x2748774cdf8eeb99",
    "0x34b0bcb5e19b48a8",
    "0x391c0cb3c5c95a63",
    "0x4ed8aa4ae3418acb",
    "0x5b9cca4f7763e373",
    "0x682e6ff3d6b2b8a3",
    "0x748f82ee5defb2fc",
    "0x78a5636f43172f60",
    "0x84c87814a1f0ab72",
    "0x8cc702081a6439ec",
    "0x90befffa23631e28",
    "0xa4506cebde82bde9",
    "0xbef9a3f7b2c67915",
    "0xc67178f2e372532b",
    "0xca273eceea26619c",
    "0xd186b8c721c0c207",
    "0xeada7dd6cde0eb1e",
    "0xf57d4f7fee6ed178",
    "0x06f067aa72176fba",
    "0x0a637dc5a2c898a6",
    "0x113f9804bef90dae",
    "0x1b710b35131c471b",
    "0x28db77f523047d84",
    "0x32caab7b40c72493",
    "0x3c9ebe0a15c9bebc",
    "0x431d67c49c100d4c",
    "0x4cc5d4becb3e42b6",
    "0x597f299cfc657e2a",
    "0x5fcb6fab3ad6faec",
    "0x6c44198c4a475817"
  ].map((n4) => BigInt(n4)));
  var SHA512_W_H = new Uint32Array(80);
  var SHA512_W_L = new Uint32Array(80);
  var SHA512 = class extends SHA2 {
    constructor() {
      super(128, 64, 16, false);
      this.Ah = 1779033703 | 0;
      this.Al = 4089235720 | 0;
      this.Bh = 3144134277 | 0;
      this.Bl = 2227873595 | 0;
      this.Ch = 1013904242 | 0;
      this.Cl = 4271175723 | 0;
      this.Dh = 2773480762 | 0;
      this.Dl = 1595750129 | 0;
      this.Eh = 1359893119 | 0;
      this.El = 2917565137 | 0;
      this.Fh = 2600822924 | 0;
      this.Fl = 725511199 | 0;
      this.Gh = 528734635 | 0;
      this.Gl = 4215389547 | 0;
      this.Hh = 1541459225 | 0;
      this.Hl = 327033209 | 0;
    }
    // prettier-ignore
    get() {
      const { Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl } = this;
      return [Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl];
    }
    // prettier-ignore
    set(Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl) {
      this.Ah = Ah | 0;
      this.Al = Al | 0;
      this.Bh = Bh | 0;
      this.Bl = Bl | 0;
      this.Ch = Ch | 0;
      this.Cl = Cl | 0;
      this.Dh = Dh | 0;
      this.Dl = Dl | 0;
      this.Eh = Eh | 0;
      this.El = El | 0;
      this.Fh = Fh | 0;
      this.Fl = Fl | 0;
      this.Gh = Gh | 0;
      this.Gl = Gl | 0;
      this.Hh = Hh | 0;
      this.Hl = Hl | 0;
    }
    process(view, offset) {
      for (let i4 = 0; i4 < 16; i4++, offset += 4) {
        SHA512_W_H[i4] = view.getUint32(offset);
        SHA512_W_L[i4] = view.getUint32(offset += 4);
      }
      for (let i4 = 16; i4 < 80; i4++) {
        const W15h = SHA512_W_H[i4 - 15] | 0;
        const W15l = SHA512_W_L[i4 - 15] | 0;
        const s0h = u64_default.rotrSH(W15h, W15l, 1) ^ u64_default.rotrSH(W15h, W15l, 8) ^ u64_default.shrSH(W15h, W15l, 7);
        const s0l = u64_default.rotrSL(W15h, W15l, 1) ^ u64_default.rotrSL(W15h, W15l, 8) ^ u64_default.shrSL(W15h, W15l, 7);
        const W2h = SHA512_W_H[i4 - 2] | 0;
        const W2l = SHA512_W_L[i4 - 2] | 0;
        const s1h = u64_default.rotrSH(W2h, W2l, 19) ^ u64_default.rotrBH(W2h, W2l, 61) ^ u64_default.shrSH(W2h, W2l, 6);
        const s1l = u64_default.rotrSL(W2h, W2l, 19) ^ u64_default.rotrBL(W2h, W2l, 61) ^ u64_default.shrSL(W2h, W2l, 6);
        const SUMl = u64_default.add4L(s0l, s1l, SHA512_W_L[i4 - 7], SHA512_W_L[i4 - 16]);
        const SUMh = u64_default.add4H(SUMl, s0h, s1h, SHA512_W_H[i4 - 7], SHA512_W_H[i4 - 16]);
        SHA512_W_H[i4] = SUMh | 0;
        SHA512_W_L[i4] = SUMl | 0;
      }
      let { Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl } = this;
      for (let i4 = 0; i4 < 80; i4++) {
        const sigma1h = u64_default.rotrSH(Eh, El, 14) ^ u64_default.rotrSH(Eh, El, 18) ^ u64_default.rotrBH(Eh, El, 41);
        const sigma1l = u64_default.rotrSL(Eh, El, 14) ^ u64_default.rotrSL(Eh, El, 18) ^ u64_default.rotrBL(Eh, El, 41);
        const CHIh = Eh & Fh ^ ~Eh & Gh;
        const CHIl = El & Fl ^ ~El & Gl;
        const T1ll = u64_default.add5L(Hl, sigma1l, CHIl, SHA512_Kl[i4], SHA512_W_L[i4]);
        const T1h = u64_default.add5H(T1ll, Hh, sigma1h, CHIh, SHA512_Kh[i4], SHA512_W_H[i4]);
        const T1l = T1ll | 0;
        const sigma0h = u64_default.rotrSH(Ah, Al, 28) ^ u64_default.rotrBH(Ah, Al, 34) ^ u64_default.rotrBH(Ah, Al, 39);
        const sigma0l = u64_default.rotrSL(Ah, Al, 28) ^ u64_default.rotrBL(Ah, Al, 34) ^ u64_default.rotrBL(Ah, Al, 39);
        const MAJh = Ah & Bh ^ Ah & Ch ^ Bh & Ch;
        const MAJl = Al & Bl ^ Al & Cl ^ Bl & Cl;
        Hh = Gh | 0;
        Hl = Gl | 0;
        Gh = Fh | 0;
        Gl = Fl | 0;
        Fh = Eh | 0;
        Fl = El | 0;
        ({ h: Eh, l: El } = u64_default.add(Dh | 0, Dl | 0, T1h | 0, T1l | 0));
        Dh = Ch | 0;
        Dl = Cl | 0;
        Ch = Bh | 0;
        Cl = Bl | 0;
        Bh = Ah | 0;
        Bl = Al | 0;
        const All = u64_default.add3L(T1l, sigma0l, MAJl);
        Ah = u64_default.add3H(All, T1h, sigma0h, MAJh);
        Al = All | 0;
      }
      ({ h: Ah, l: Al } = u64_default.add(this.Ah | 0, this.Al | 0, Ah | 0, Al | 0));
      ({ h: Bh, l: Bl } = u64_default.add(this.Bh | 0, this.Bl | 0, Bh | 0, Bl | 0));
      ({ h: Ch, l: Cl } = u64_default.add(this.Ch | 0, this.Cl | 0, Ch | 0, Cl | 0));
      ({ h: Dh, l: Dl } = u64_default.add(this.Dh | 0, this.Dl | 0, Dh | 0, Dl | 0));
      ({ h: Eh, l: El } = u64_default.add(this.Eh | 0, this.El | 0, Eh | 0, El | 0));
      ({ h: Fh, l: Fl } = u64_default.add(this.Fh | 0, this.Fl | 0, Fh | 0, Fl | 0));
      ({ h: Gh, l: Gl } = u64_default.add(this.Gh | 0, this.Gl | 0, Gh | 0, Gl | 0));
      ({ h: Hh, l: Hl } = u64_default.add(this.Hh | 0, this.Hl | 0, Hh | 0, Hl | 0));
      this.set(Ah, Al, Bh, Bl, Ch, Cl, Dh, Dl, Eh, El, Fh, Fl, Gh, Gl, Hh, Hl);
    }
    roundClean() {
      SHA512_W_H.fill(0);
      SHA512_W_L.fill(0);
    }
    destroy() {
      this.buffer.fill(0);
      this.set(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0);
    }
  };
  var SHA512_256 = class extends SHA512 {
    constructor() {
      super();
      this.Ah = 573645204 | 0;
      this.Al = 4230739756 | 0;
      this.Bh = 2673172387 | 0;
      this.Bl = 3360449730 | 0;
      this.Ch = 596883563 | 0;
      this.Cl = 1867755857 | 0;
      this.Dh = 2520282905 | 0;
      this.Dl = 1497426621 | 0;
      this.Eh = 2519219938 | 0;
      this.El = 2827943907 | 0;
      this.Fh = 3193839141 | 0;
      this.Fl = 1401305490 | 0;
      this.Gh = 721525244 | 0;
      this.Gl = 746961066 | 0;
      this.Hh = 246885852 | 0;
      this.Hl = 2177182882 | 0;
      this.outputLen = 32;
    }
  };
  var SHA384 = class extends SHA512 {
    constructor() {
      super();
      this.Ah = 3418070365 | 0;
      this.Al = 3238371032 | 0;
      this.Bh = 1654270250 | 0;
      this.Bl = 914150663 | 0;
      this.Ch = 2438529370 | 0;
      this.Cl = 812702999 | 0;
      this.Dh = 355462360 | 0;
      this.Dl = 4144912697 | 0;
      this.Eh = 1731405415 | 0;
      this.El = 4290775857 | 0;
      this.Fh = 2394180231 | 0;
      this.Fl = 1750603025 | 0;
      this.Gh = 3675008525 | 0;
      this.Gl = 1694076839 | 0;
      this.Hh = 1203062813 | 0;
      this.Hl = 3204075428 | 0;
      this.outputLen = 48;
    }
  };
  var sha512 = wrapConstructor(() => new SHA512());
  var sha512_256 = wrapConstructor(() => new SHA512_256());
  var sha384 = wrapConstructor(() => new SHA384());

  // node_modules/ethers/lib.esm/crypto/crypto-browser.js
  function getGlobal() {
    if (typeof self !== "undefined") {
      return self;
    }
    if (typeof window !== "undefined") {
      return window;
    }
    if (typeof global !== "undefined") {
      return global;
    }
    throw new Error("unable to locate global object");
  }
  var anyGlobal = getGlobal();
  var crypto2 = anyGlobal.crypto || anyGlobal.msCrypto;
  function createHmac(_algo, key) {
    const algo = { sha256, sha512 }[_algo];
    assertArgument(algo != null, "invalid hmac algorithm", "algorithm", _algo);
    return hmac.create(algo, key);
  }

  // node_modules/ethers/lib.esm/crypto/hmac.js
  var locked2 = false;
  var _computeHmac = function(algorithm, key, data) {
    return createHmac(algorithm, key).update(data).digest();
  };
  var __computeHmac = _computeHmac;
  function computeHmac(algorithm, _key, _data) {
    const key = getBytes(_key, "key");
    const data = getBytes(_data, "data");
    return hexlify(__computeHmac(algorithm, key, data));
  }
  computeHmac._ = _computeHmac;
  computeHmac.lock = function() {
    locked2 = true;
  };
  computeHmac.register = function(func) {
    if (locked2) {
      throw new Error("computeHmac is locked");
    }
    __computeHmac = func;
  };
  Object.freeze(computeHmac);

  // node_modules/@noble/hashes/esm/sha3.js
  var [SHA3_PI, SHA3_ROTL, _SHA3_IOTA] = [[], [], []];
  var _0n = BigInt(0);
  var _1n = BigInt(1);
  var _2n = BigInt(2);
  var _7n = BigInt(7);
  var _256n = BigInt(256);
  var _0x71n = BigInt(113);
  for (let round = 0, R = _1n, x2 = 1, y2 = 0; round < 24; round++) {
    [x2, y2] = [y2, (2 * x2 + 3 * y2) % 5];
    SHA3_PI.push(2 * (5 * y2 + x2));
    SHA3_ROTL.push((round + 1) * (round + 2) / 2 % 64);
    let t4 = _0n;
    for (let j4 = 0; j4 < 7; j4++) {
      R = (R << _1n ^ (R >> _7n) * _0x71n) % _256n;
      if (R & _2n)
        t4 ^= _1n << (_1n << BigInt(j4)) - _1n;
    }
    _SHA3_IOTA.push(t4);
  }
  var [SHA3_IOTA_H, SHA3_IOTA_L] = u64_default.split(_SHA3_IOTA, true);
  var rotlH = (h4, l4, s4) => s4 > 32 ? u64_default.rotlBH(h4, l4, s4) : u64_default.rotlSH(h4, l4, s4);
  var rotlL = (h4, l4, s4) => s4 > 32 ? u64_default.rotlBL(h4, l4, s4) : u64_default.rotlSL(h4, l4, s4);
  function keccakP(s4, rounds = 24) {
    const B3 = new Uint32Array(5 * 2);
    for (let round = 24 - rounds; round < 24; round++) {
      for (let x2 = 0; x2 < 10; x2++)
        B3[x2] = s4[x2] ^ s4[x2 + 10] ^ s4[x2 + 20] ^ s4[x2 + 30] ^ s4[x2 + 40];
      for (let x2 = 0; x2 < 10; x2 += 2) {
        const idx1 = (x2 + 8) % 10;
        const idx0 = (x2 + 2) % 10;
        const B0 = B3[idx0];
        const B1 = B3[idx0 + 1];
        const Th = rotlH(B0, B1, 1) ^ B3[idx1];
        const Tl = rotlL(B0, B1, 1) ^ B3[idx1 + 1];
        for (let y2 = 0; y2 < 50; y2 += 10) {
          s4[x2 + y2] ^= Th;
          s4[x2 + y2 + 1] ^= Tl;
        }
      }
      let curH = s4[2];
      let curL = s4[3];
      for (let t4 = 0; t4 < 24; t4++) {
        const shift = SHA3_ROTL[t4];
        const Th = rotlH(curH, curL, shift);
        const Tl = rotlL(curH, curL, shift);
        const PI = SHA3_PI[t4];
        curH = s4[PI];
        curL = s4[PI + 1];
        s4[PI] = Th;
        s4[PI + 1] = Tl;
      }
      for (let y2 = 0; y2 < 50; y2 += 10) {
        for (let x2 = 0; x2 < 10; x2++)
          B3[x2] = s4[y2 + x2];
        for (let x2 = 0; x2 < 10; x2++)
          s4[y2 + x2] ^= ~B3[(x2 + 2) % 10] & B3[(x2 + 4) % 10];
      }
      s4[0] ^= SHA3_IOTA_H[round];
      s4[1] ^= SHA3_IOTA_L[round];
    }
    B3.fill(0);
  }
  var Keccak = class _Keccak extends Hash {
    // NOTE: we accept arguments in bytes instead of bits here.
    constructor(blockLen, suffix, outputLen, enableXOF = false, rounds = 24) {
      super();
      this.blockLen = blockLen;
      this.suffix = suffix;
      this.outputLen = outputLen;
      this.enableXOF = enableXOF;
      this.rounds = rounds;
      this.pos = 0;
      this.posOut = 0;
      this.finished = false;
      this.destroyed = false;
      assert_default.number(outputLen);
      if (0 >= this.blockLen || this.blockLen >= 200)
        throw new Error("Sha3 supports only keccak-f1600 function");
      this.state = new Uint8Array(200);
      this.state32 = u32(this.state);
    }
    keccak() {
      keccakP(this.state32, this.rounds);
      this.posOut = 0;
      this.pos = 0;
    }
    update(data) {
      assert_default.exists(this);
      const { blockLen, state } = this;
      data = toBytes(data);
      const len = data.length;
      for (let pos = 0; pos < len; ) {
        const take = Math.min(blockLen - this.pos, len - pos);
        for (let i4 = 0; i4 < take; i4++)
          state[this.pos++] ^= data[pos++];
        if (this.pos === blockLen)
          this.keccak();
      }
      return this;
    }
    finish() {
      if (this.finished)
        return;
      this.finished = true;
      const { state, suffix, pos, blockLen } = this;
      state[pos] ^= suffix;
      if ((suffix & 128) !== 0 && pos === blockLen - 1)
        this.keccak();
      state[blockLen - 1] ^= 128;
      this.keccak();
    }
    writeInto(out) {
      assert_default.exists(this, false);
      assert_default.bytes(out);
      this.finish();
      const bufferOut = this.state;
      const { blockLen } = this;
      for (let pos = 0, len = out.length; pos < len; ) {
        if (this.posOut >= blockLen)
          this.keccak();
        const take = Math.min(blockLen - this.posOut, len - pos);
        out.set(bufferOut.subarray(this.posOut, this.posOut + take), pos);
        this.posOut += take;
        pos += take;
      }
      return out;
    }
    xofInto(out) {
      if (!this.enableXOF)
        throw new Error("XOF is not possible for this instance");
      return this.writeInto(out);
    }
    xof(bytes2) {
      assert_default.number(bytes2);
      return this.xofInto(new Uint8Array(bytes2));
    }
    digestInto(out) {
      assert_default.output(out, this);
      if (this.finished)
        throw new Error("digest() was already called");
      this.writeInto(out);
      this.destroy();
      return out;
    }
    digest() {
      return this.digestInto(new Uint8Array(this.outputLen));
    }
    destroy() {
      this.destroyed = true;
      this.state.fill(0);
    }
    _cloneInto(to) {
      const { blockLen, suffix, outputLen, rounds, enableXOF } = this;
      to || (to = new _Keccak(blockLen, suffix, outputLen, enableXOF, rounds));
      to.state32.set(this.state32);
      to.pos = this.pos;
      to.posOut = this.posOut;
      to.finished = this.finished;
      to.rounds = rounds;
      to.suffix = suffix;
      to.outputLen = outputLen;
      to.enableXOF = enableXOF;
      to.destroyed = this.destroyed;
      return to;
    }
  };
  var gen = (suffix, blockLen, outputLen) => wrapConstructor(() => new Keccak(blockLen, suffix, outputLen));
  var sha3_224 = gen(6, 144, 224 / 8);
  var sha3_256 = gen(6, 136, 256 / 8);
  var sha3_384 = gen(6, 104, 384 / 8);
  var sha3_512 = gen(6, 72, 512 / 8);
  var keccak_224 = gen(1, 144, 224 / 8);
  var keccak_256 = gen(1, 136, 256 / 8);
  var keccak_384 = gen(1, 104, 384 / 8);
  var keccak_512 = gen(1, 72, 512 / 8);
  var genShake = (suffix, blockLen, outputLen) => wrapConstructorWithOpts((opts = {}) => new Keccak(blockLen, suffix, opts.dkLen === void 0 ? outputLen : opts.dkLen, true));
  var shake128 = genShake(31, 168, 128 / 8);
  var shake256 = genShake(31, 136, 256 / 8);

  // node_modules/ethers/lib.esm/crypto/keccak.js
  var locked3 = false;
  var _keccak256 = function(data) {
    return keccak_256(data);
  };
  var __keccak256 = _keccak256;
  function keccak256(_data) {
    const data = getBytes(_data, "data");
    return hexlify(__keccak256(data));
  }
  keccak256._ = _keccak256;
  keccak256.lock = function() {
    locked3 = true;
  };
  keccak256.register = function(func) {
    if (locked3) {
      throw new TypeError("keccak256 is locked");
    }
    __keccak256 = func;
  };
  Object.freeze(keccak256);

  // node_modules/@noble/secp256k1/lib/esm/index.js
  var nodeCrypto = __toESM(require_crypto(), 1);
  var _0n2 = BigInt(0);
  var _1n2 = BigInt(1);
  var _2n2 = BigInt(2);
  var _3n = BigInt(3);
  var _8n = BigInt(8);
  var CURVE = Object.freeze({
    a: _0n2,
    b: BigInt(7),
    P: BigInt("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffffc2f"),
    n: BigInt("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
    h: _1n2,
    Gx: BigInt("55066263022277343669578718895168534326250603453777594175500187360389116729240"),
    Gy: BigInt("32670510020758816978083085130507043184471273380659243275938904335757337482424"),
    beta: BigInt("0x7ae96a2b657c07106e64479eac3434e99cf0497512f58995c1396c28719501ee")
  });
  var divNearest = (a4, b5) => (a4 + b5 / _2n2) / b5;
  var endo = {
    beta: BigInt("0x7ae96a2b657c07106e64479eac3434e99cf0497512f58995c1396c28719501ee"),
    splitScalar(k3) {
      const { n: n4 } = CURVE;
      const a1 = BigInt("0x3086d221a7d46bcde86c90e49284eb15");
      const b1 = -_1n2 * BigInt("0xe4437ed6010e88286f547fa90abfe4c3");
      const a22 = BigInt("0x114ca50f7a8e2f3f657c1108d9d44cfd8");
      const b22 = a1;
      const POW_2_128 = BigInt("0x100000000000000000000000000000000");
      const c1 = divNearest(b22 * k3, n4);
      const c22 = divNearest(-b1 * k3, n4);
      let k1 = mod(k3 - c1 * a1 - c22 * a22, n4);
      let k22 = mod(-c1 * b1 - c22 * b22, n4);
      const k1neg = k1 > POW_2_128;
      const k2neg = k22 > POW_2_128;
      if (k1neg)
        k1 = n4 - k1;
      if (k2neg)
        k22 = n4 - k22;
      if (k1 > POW_2_128 || k22 > POW_2_128) {
        throw new Error("splitScalarEndo: Endomorphism failed, k=" + k3);
      }
      return { k1neg, k1, k2neg, k2: k22 };
    }
  };
  var fieldLen = 32;
  var groupLen = 32;
  var hashLen = 32;
  var compressedLen = fieldLen + 1;
  var uncompressedLen = 2 * fieldLen + 1;
  function weierstrass(x2) {
    const { a: a4, b: b5 } = CURVE;
    const x22 = mod(x2 * x2);
    const x3 = mod(x22 * x2);
    return mod(x3 + a4 * x2 + b5);
  }
  var USE_ENDOMORPHISM = CURVE.a === _0n2;
  var ShaError = class extends Error {
    constructor(message) {
      super(message);
    }
  };
  function assertJacPoint(other) {
    if (!(other instanceof JacobianPoint))
      throw new TypeError("JacobianPoint expected");
  }
  var JacobianPoint = class _JacobianPoint {
    constructor(x2, y2, z3) {
      this.x = x2;
      this.y = y2;
      this.z = z3;
    }
    static fromAffine(p4) {
      if (!(p4 instanceof Point)) {
        throw new TypeError("JacobianPoint#fromAffine: expected Point");
      }
      if (p4.equals(Point.ZERO))
        return _JacobianPoint.ZERO;
      return new _JacobianPoint(p4.x, p4.y, _1n2);
    }
    static toAffineBatch(points) {
      const toInv = invertBatch(points.map((p4) => p4.z));
      return points.map((p4, i4) => p4.toAffine(toInv[i4]));
    }
    static normalizeZ(points) {
      return _JacobianPoint.toAffineBatch(points).map(_JacobianPoint.fromAffine);
    }
    equals(other) {
      assertJacPoint(other);
      const { x: X1, y: Y1, z: Z1 } = this;
      const { x: X2, y: Y2, z: Z2 } = other;
      const Z1Z1 = mod(Z1 * Z1);
      const Z2Z2 = mod(Z2 * Z2);
      const U1 = mod(X1 * Z2Z2);
      const U2 = mod(X2 * Z1Z1);
      const S12 = mod(mod(Y1 * Z2) * Z2Z2);
      const S2 = mod(mod(Y2 * Z1) * Z1Z1);
      return U1 === U2 && S12 === S2;
    }
    negate() {
      return new _JacobianPoint(this.x, mod(-this.y), this.z);
    }
    double() {
      const { x: X1, y: Y1, z: Z1 } = this;
      const A2 = mod(X1 * X1);
      const B3 = mod(Y1 * Y1);
      const C = mod(B3 * B3);
      const x1b = X1 + B3;
      const D2 = mod(_2n2 * (mod(x1b * x1b) - A2 - C));
      const E = mod(_3n * A2);
      const F2 = mod(E * E);
      const X3 = mod(F2 - _2n2 * D2);
      const Y3 = mod(E * (D2 - X3) - _8n * C);
      const Z3 = mod(_2n2 * Y1 * Z1);
      return new _JacobianPoint(X3, Y3, Z3);
    }
    add(other) {
      assertJacPoint(other);
      const { x: X1, y: Y1, z: Z1 } = this;
      const { x: X2, y: Y2, z: Z2 } = other;
      if (X2 === _0n2 || Y2 === _0n2)
        return this;
      if (X1 === _0n2 || Y1 === _0n2)
        return other;
      const Z1Z1 = mod(Z1 * Z1);
      const Z2Z2 = mod(Z2 * Z2);
      const U1 = mod(X1 * Z2Z2);
      const U2 = mod(X2 * Z1Z1);
      const S12 = mod(mod(Y1 * Z2) * Z2Z2);
      const S2 = mod(mod(Y2 * Z1) * Z1Z1);
      const H2 = mod(U2 - U1);
      const r4 = mod(S2 - S12);
      if (H2 === _0n2) {
        if (r4 === _0n2) {
          return this.double();
        } else {
          return _JacobianPoint.ZERO;
        }
      }
      const HH = mod(H2 * H2);
      const HHH = mod(H2 * HH);
      const V = mod(U1 * HH);
      const X3 = mod(r4 * r4 - HHH - _2n2 * V);
      const Y3 = mod(r4 * (V - X3) - S12 * HHH);
      const Z3 = mod(Z1 * Z2 * H2);
      return new _JacobianPoint(X3, Y3, Z3);
    }
    subtract(other) {
      return this.add(other.negate());
    }
    multiplyUnsafe(scalar) {
      const P0 = _JacobianPoint.ZERO;
      if (typeof scalar === "bigint" && scalar === _0n2)
        return P0;
      let n4 = normalizeScalar(scalar);
      if (n4 === _1n2)
        return this;
      if (!USE_ENDOMORPHISM) {
        let p4 = P0;
        let d5 = this;
        while (n4 > _0n2) {
          if (n4 & _1n2)
            p4 = p4.add(d5);
          d5 = d5.double();
          n4 >>= _1n2;
        }
        return p4;
      }
      let { k1neg, k1, k2neg, k2: k22 } = endo.splitScalar(n4);
      let k1p = P0;
      let k2p = P0;
      let d4 = this;
      while (k1 > _0n2 || k22 > _0n2) {
        if (k1 & _1n2)
          k1p = k1p.add(d4);
        if (k22 & _1n2)
          k2p = k2p.add(d4);
        d4 = d4.double();
        k1 >>= _1n2;
        k22 >>= _1n2;
      }
      if (k1neg)
        k1p = k1p.negate();
      if (k2neg)
        k2p = k2p.negate();
      k2p = new _JacobianPoint(mod(k2p.x * endo.beta), k2p.y, k2p.z);
      return k1p.add(k2p);
    }
    precomputeWindow(W) {
      const windows = USE_ENDOMORPHISM ? 128 / W + 1 : 256 / W + 1;
      const points = [];
      let p4 = this;
      let base = p4;
      for (let window2 = 0; window2 < windows; window2++) {
        base = p4;
        points.push(base);
        for (let i4 = 1; i4 < 2 ** (W - 1); i4++) {
          base = base.add(p4);
          points.push(base);
        }
        p4 = base.double();
      }
      return points;
    }
    wNAF(n4, affinePoint) {
      if (!affinePoint && this.equals(_JacobianPoint.BASE))
        affinePoint = Point.BASE;
      const W = affinePoint && affinePoint._WINDOW_SIZE || 1;
      if (256 % W) {
        throw new Error("Point#wNAF: Invalid precomputation window, must be power of 2");
      }
      let precomputes = affinePoint && pointPrecomputes.get(affinePoint);
      if (!precomputes) {
        precomputes = this.precomputeWindow(W);
        if (affinePoint && W !== 1) {
          precomputes = _JacobianPoint.normalizeZ(precomputes);
          pointPrecomputes.set(affinePoint, precomputes);
        }
      }
      let p4 = _JacobianPoint.ZERO;
      let f4 = _JacobianPoint.BASE;
      const windows = 1 + (USE_ENDOMORPHISM ? 128 / W : 256 / W);
      const windowSize = 2 ** (W - 1);
      const mask2 = BigInt(2 ** W - 1);
      const maxNumber = 2 ** W;
      const shiftBy = BigInt(W);
      for (let window2 = 0; window2 < windows; window2++) {
        const offset = window2 * windowSize;
        let wbits = Number(n4 & mask2);
        n4 >>= shiftBy;
        if (wbits > windowSize) {
          wbits -= maxNumber;
          n4 += _1n2;
        }
        const offset1 = offset;
        const offset2 = offset + Math.abs(wbits) - 1;
        const cond1 = window2 % 2 !== 0;
        const cond2 = wbits < 0;
        if (wbits === 0) {
          f4 = f4.add(constTimeNegate(cond1, precomputes[offset1]));
        } else {
          p4 = p4.add(constTimeNegate(cond2, precomputes[offset2]));
        }
      }
      return { p: p4, f: f4 };
    }
    multiply(scalar, affinePoint) {
      let n4 = normalizeScalar(scalar);
      let point;
      let fake;
      if (USE_ENDOMORPHISM) {
        const { k1neg, k1, k2neg, k2: k22 } = endo.splitScalar(n4);
        let { p: k1p, f: f1p } = this.wNAF(k1, affinePoint);
        let { p: k2p, f: f2p } = this.wNAF(k22, affinePoint);
        k1p = constTimeNegate(k1neg, k1p);
        k2p = constTimeNegate(k2neg, k2p);
        k2p = new _JacobianPoint(mod(k2p.x * endo.beta), k2p.y, k2p.z);
        point = k1p.add(k2p);
        fake = f1p.add(f2p);
      } else {
        const { p: p4, f: f4 } = this.wNAF(n4, affinePoint);
        point = p4;
        fake = f4;
      }
      return _JacobianPoint.normalizeZ([point, fake])[0];
    }
    toAffine(invZ) {
      const { x: x2, y: y2, z: z3 } = this;
      const is0 = this.equals(_JacobianPoint.ZERO);
      if (invZ == null)
        invZ = is0 ? _8n : invert(z3);
      const iz1 = invZ;
      const iz2 = mod(iz1 * iz1);
      const iz3 = mod(iz2 * iz1);
      const ax = mod(x2 * iz2);
      const ay = mod(y2 * iz3);
      const zz = mod(z3 * iz1);
      if (is0)
        return Point.ZERO;
      if (zz !== _1n2)
        throw new Error("invZ was invalid");
      return new Point(ax, ay);
    }
  };
  JacobianPoint.BASE = new JacobianPoint(CURVE.Gx, CURVE.Gy, _1n2);
  JacobianPoint.ZERO = new JacobianPoint(_0n2, _1n2, _0n2);
  function constTimeNegate(condition, item) {
    const neg = item.negate();
    return condition ? neg : item;
  }
  var pointPrecomputes = /* @__PURE__ */ new WeakMap();
  var Point = class _Point {
    constructor(x2, y2) {
      this.x = x2;
      this.y = y2;
    }
    _setWindowSize(windowSize) {
      this._WINDOW_SIZE = windowSize;
      pointPrecomputes.delete(this);
    }
    hasEvenY() {
      return this.y % _2n2 === _0n2;
    }
    static fromCompressedHex(bytes2) {
      const isShort = bytes2.length === 32;
      const x2 = bytesToNumber(isShort ? bytes2 : bytes2.subarray(1));
      if (!isValidFieldElement(x2))
        throw new Error("Point is not on curve");
      const y2 = weierstrass(x2);
      let y3 = sqrtMod(y2);
      const isYOdd = (y3 & _1n2) === _1n2;
      if (isShort) {
        if (isYOdd)
          y3 = mod(-y3);
      } else {
        const isFirstByteOdd = (bytes2[0] & 1) === 1;
        if (isFirstByteOdd !== isYOdd)
          y3 = mod(-y3);
      }
      const point = new _Point(x2, y3);
      point.assertValidity();
      return point;
    }
    static fromUncompressedHex(bytes2) {
      const x2 = bytesToNumber(bytes2.subarray(1, fieldLen + 1));
      const y2 = bytesToNumber(bytes2.subarray(fieldLen + 1, fieldLen * 2 + 1));
      const point = new _Point(x2, y2);
      point.assertValidity();
      return point;
    }
    static fromHex(hex) {
      const bytes2 = ensureBytes(hex);
      const len = bytes2.length;
      const header = bytes2[0];
      if (len === fieldLen)
        return this.fromCompressedHex(bytes2);
      if (len === compressedLen && (header === 2 || header === 3)) {
        return this.fromCompressedHex(bytes2);
      }
      if (len === uncompressedLen && header === 4)
        return this.fromUncompressedHex(bytes2);
      throw new Error(`Point.fromHex: received invalid point. Expected 32-${compressedLen} compressed bytes or ${uncompressedLen} uncompressed bytes, not ${len}`);
    }
    static fromPrivateKey(privateKey) {
      return _Point.BASE.multiply(normalizePrivateKey(privateKey));
    }
    static fromSignature(msgHash, signature, recovery) {
      const { r: r4, s: s4 } = normalizeSignature(signature);
      if (![0, 1, 2, 3].includes(recovery))
        throw new Error("Cannot recover: invalid recovery bit");
      const h4 = truncateHash(ensureBytes(msgHash));
      const { n: n4 } = CURVE;
      const radj = recovery === 2 || recovery === 3 ? r4 + n4 : r4;
      const rinv = invert(radj, n4);
      const u1 = mod(-h4 * rinv, n4);
      const u22 = mod(s4 * rinv, n4);
      const prefix = recovery & 1 ? "03" : "02";
      const R = _Point.fromHex(prefix + numTo32bStr(radj));
      const Q = _Point.BASE.multiplyAndAddUnsafe(R, u1, u22);
      if (!Q)
        throw new Error("Cannot recover signature: point at infinify");
      Q.assertValidity();
      return Q;
    }
    toRawBytes(isCompressed = false) {
      return hexToBytes(this.toHex(isCompressed));
    }
    toHex(isCompressed = false) {
      const x2 = numTo32bStr(this.x);
      if (isCompressed) {
        const prefix = this.hasEvenY() ? "02" : "03";
        return `${prefix}${x2}`;
      } else {
        return `04${x2}${numTo32bStr(this.y)}`;
      }
    }
    toHexX() {
      return this.toHex(true).slice(2);
    }
    toRawX() {
      return this.toRawBytes(true).slice(1);
    }
    assertValidity() {
      const msg = "Point is not on elliptic curve";
      const { x: x2, y: y2 } = this;
      if (!isValidFieldElement(x2) || !isValidFieldElement(y2))
        throw new Error(msg);
      const left = mod(y2 * y2);
      const right = weierstrass(x2);
      if (mod(left - right) !== _0n2)
        throw new Error(msg);
    }
    equals(other) {
      return this.x === other.x && this.y === other.y;
    }
    negate() {
      return new _Point(this.x, mod(-this.y));
    }
    double() {
      return JacobianPoint.fromAffine(this).double().toAffine();
    }
    add(other) {
      return JacobianPoint.fromAffine(this).add(JacobianPoint.fromAffine(other)).toAffine();
    }
    subtract(other) {
      return this.add(other.negate());
    }
    multiply(scalar) {
      return JacobianPoint.fromAffine(this).multiply(scalar, this).toAffine();
    }
    multiplyAndAddUnsafe(Q, a4, b5) {
      const P2 = JacobianPoint.fromAffine(this);
      const aP = a4 === _0n2 || a4 === _1n2 || this !== _Point.BASE ? P2.multiplyUnsafe(a4) : P2.multiply(a4);
      const bQ = JacobianPoint.fromAffine(Q).multiplyUnsafe(b5);
      const sum = aP.add(bQ);
      return sum.equals(JacobianPoint.ZERO) ? void 0 : sum.toAffine();
    }
  };
  Point.BASE = new Point(CURVE.Gx, CURVE.Gy);
  Point.ZERO = new Point(_0n2, _0n2);
  function sliceDER(s4) {
    return Number.parseInt(s4[0], 16) >= 8 ? "00" + s4 : s4;
  }
  function parseDERInt(data) {
    if (data.length < 2 || data[0] !== 2) {
      throw new Error(`Invalid signature integer tag: ${bytesToHex(data)}`);
    }
    const len = data[1];
    const res = data.subarray(2, len + 2);
    if (!len || res.length !== len) {
      throw new Error(`Invalid signature integer: wrong length`);
    }
    if (res[0] === 0 && res[1] <= 127) {
      throw new Error("Invalid signature integer: trailing length");
    }
    return { data: bytesToNumber(res), left: data.subarray(len + 2) };
  }
  function parseDERSignature(data) {
    if (data.length < 2 || data[0] != 48) {
      throw new Error(`Invalid signature tag: ${bytesToHex(data)}`);
    }
    if (data[1] !== data.length - 2) {
      throw new Error("Invalid signature: incorrect length");
    }
    const { data: r4, left: sBytes } = parseDERInt(data.subarray(2));
    const { data: s4, left: rBytesLeft } = parseDERInt(sBytes);
    if (rBytesLeft.length) {
      throw new Error(`Invalid signature: left bytes after parsing: ${bytesToHex(rBytesLeft)}`);
    }
    return { r: r4, s: s4 };
  }
  var Signature = class _Signature {
    constructor(r4, s4) {
      this.r = r4;
      this.s = s4;
      this.assertValidity();
    }
    static fromCompact(hex) {
      const arr = hex instanceof Uint8Array;
      const name = "Signature.fromCompact";
      if (typeof hex !== "string" && !arr)
        throw new TypeError(`${name}: Expected string or Uint8Array`);
      const str = arr ? bytesToHex(hex) : hex;
      if (str.length !== 128)
        throw new Error(`${name}: Expected 64-byte hex`);
      return new _Signature(hexToNumber(str.slice(0, 64)), hexToNumber(str.slice(64, 128)));
    }
    static fromDER(hex) {
      const arr = hex instanceof Uint8Array;
      if (typeof hex !== "string" && !arr)
        throw new TypeError(`Signature.fromDER: Expected string or Uint8Array`);
      const { r: r4, s: s4 } = parseDERSignature(arr ? hex : hexToBytes(hex));
      return new _Signature(r4, s4);
    }
    static fromHex(hex) {
      return this.fromDER(hex);
    }
    assertValidity() {
      const { r: r4, s: s4 } = this;
      if (!isWithinCurveOrder(r4))
        throw new Error("Invalid Signature: r must be 0 < r < n");
      if (!isWithinCurveOrder(s4))
        throw new Error("Invalid Signature: s must be 0 < s < n");
    }
    hasHighS() {
      const HALF = CURVE.n >> _1n2;
      return this.s > HALF;
    }
    normalizeS() {
      return this.hasHighS() ? new _Signature(this.r, mod(-this.s, CURVE.n)) : this;
    }
    toDERRawBytes() {
      return hexToBytes(this.toDERHex());
    }
    toDERHex() {
      const sHex = sliceDER(numberToHexUnpadded(this.s));
      const rHex = sliceDER(numberToHexUnpadded(this.r));
      const sHexL = sHex.length / 2;
      const rHexL = rHex.length / 2;
      const sLen = numberToHexUnpadded(sHexL);
      const rLen = numberToHexUnpadded(rHexL);
      const length = numberToHexUnpadded(rHexL + sHexL + 4);
      return `30${length}02${rLen}${rHex}02${sLen}${sHex}`;
    }
    toRawBytes() {
      return this.toDERRawBytes();
    }
    toHex() {
      return this.toDERHex();
    }
    toCompactRawBytes() {
      return hexToBytes(this.toCompactHex());
    }
    toCompactHex() {
      return numTo32bStr(this.r) + numTo32bStr(this.s);
    }
  };
  function concatBytes(...arrays) {
    if (!arrays.every((b5) => b5 instanceof Uint8Array))
      throw new Error("Uint8Array list expected");
    if (arrays.length === 1)
      return arrays[0];
    const length = arrays.reduce((a4, arr) => a4 + arr.length, 0);
    const result = new Uint8Array(length);
    for (let i4 = 0, pad = 0; i4 < arrays.length; i4++) {
      const arr = arrays[i4];
      result.set(arr, pad);
      pad += arr.length;
    }
    return result;
  }
  var hexes2 = Array.from({ length: 256 }, (v3, i4) => i4.toString(16).padStart(2, "0"));
  function bytesToHex(uint8a) {
    if (!(uint8a instanceof Uint8Array))
      throw new Error("Expected Uint8Array");
    let hex = "";
    for (let i4 = 0; i4 < uint8a.length; i4++) {
      hex += hexes2[uint8a[i4]];
    }
    return hex;
  }
  var POW_2_256 = BigInt("0x10000000000000000000000000000000000000000000000000000000000000000");
  function numTo32bStr(num) {
    if (typeof num !== "bigint")
      throw new Error("Expected bigint");
    if (!(_0n2 <= num && num < POW_2_256))
      throw new Error("Expected number 0 <= n < 2^256");
    return num.toString(16).padStart(64, "0");
  }
  function numTo32b(num) {
    const b5 = hexToBytes(numTo32bStr(num));
    if (b5.length !== 32)
      throw new Error("Error: expected 32 bytes");
    return b5;
  }
  function numberToHexUnpadded(num) {
    const hex = num.toString(16);
    return hex.length & 1 ? `0${hex}` : hex;
  }
  function hexToNumber(hex) {
    if (typeof hex !== "string") {
      throw new TypeError("hexToNumber: expected string, got " + typeof hex);
    }
    return BigInt(`0x${hex}`);
  }
  function hexToBytes(hex) {
    if (typeof hex !== "string") {
      throw new TypeError("hexToBytes: expected string, got " + typeof hex);
    }
    if (hex.length % 2)
      throw new Error("hexToBytes: received invalid unpadded hex" + hex.length);
    const array = new Uint8Array(hex.length / 2);
    for (let i4 = 0; i4 < array.length; i4++) {
      const j4 = i4 * 2;
      const hexByte = hex.slice(j4, j4 + 2);
      const byte = Number.parseInt(hexByte, 16);
      if (Number.isNaN(byte) || byte < 0)
        throw new Error("Invalid byte sequence");
      array[i4] = byte;
    }
    return array;
  }
  function bytesToNumber(bytes2) {
    return hexToNumber(bytesToHex(bytes2));
  }
  function ensureBytes(hex) {
    return hex instanceof Uint8Array ? Uint8Array.from(hex) : hexToBytes(hex);
  }
  function normalizeScalar(num) {
    if (typeof num === "number" && Number.isSafeInteger(num) && num > 0)
      return BigInt(num);
    if (typeof num === "bigint" && isWithinCurveOrder(num))
      return num;
    throw new TypeError("Expected valid private scalar: 0 < scalar < curve.n");
  }
  function mod(a4, b5 = CURVE.P) {
    const result = a4 % b5;
    return result >= _0n2 ? result : b5 + result;
  }
  function pow2(x2, power) {
    const { P: P2 } = CURVE;
    let res = x2;
    while (power-- > _0n2) {
      res *= res;
      res %= P2;
    }
    return res;
  }
  function sqrtMod(x2) {
    const { P: P2 } = CURVE;
    const _6n = BigInt(6);
    const _11n = BigInt(11);
    const _22n = BigInt(22);
    const _23n = BigInt(23);
    const _44n = BigInt(44);
    const _88n = BigInt(88);
    const b22 = x2 * x2 * x2 % P2;
    const b32 = b22 * b22 * x2 % P2;
    const b6 = pow2(b32, _3n) * b32 % P2;
    const b9 = pow2(b6, _3n) * b32 % P2;
    const b11 = pow2(b9, _2n2) * b22 % P2;
    const b222 = pow2(b11, _11n) * b11 % P2;
    const b44 = pow2(b222, _22n) * b222 % P2;
    const b88 = pow2(b44, _44n) * b44 % P2;
    const b176 = pow2(b88, _88n) * b88 % P2;
    const b220 = pow2(b176, _44n) * b44 % P2;
    const b223 = pow2(b220, _3n) * b32 % P2;
    const t1 = pow2(b223, _23n) * b222 % P2;
    const t22 = pow2(t1, _6n) * b22 % P2;
    const rt = pow2(t22, _2n2);
    const xc = rt * rt % P2;
    if (xc !== x2)
      throw new Error("Cannot find square root");
    return rt;
  }
  function invert(number2, modulo = CURVE.P) {
    if (number2 === _0n2 || modulo <= _0n2) {
      throw new Error(`invert: expected positive integers, got n=${number2} mod=${modulo}`);
    }
    let a4 = mod(number2, modulo);
    let b5 = modulo;
    let x2 = _0n2, y2 = _1n2, u4 = _1n2, v3 = _0n2;
    while (a4 !== _0n2) {
      const q2 = b5 / a4;
      const r4 = b5 % a4;
      const m4 = x2 - u4 * q2;
      const n4 = y2 - v3 * q2;
      b5 = a4, a4 = r4, x2 = u4, y2 = v3, u4 = m4, v3 = n4;
    }
    const gcd = b5;
    if (gcd !== _1n2)
      throw new Error("invert: does not exist");
    return mod(x2, modulo);
  }
  function invertBatch(nums, p4 = CURVE.P) {
    const scratch = new Array(nums.length);
    const lastMultiplied = nums.reduce((acc, num, i4) => {
      if (num === _0n2)
        return acc;
      scratch[i4] = acc;
      return mod(acc * num, p4);
    }, _1n2);
    const inverted = invert(lastMultiplied, p4);
    nums.reduceRight((acc, num, i4) => {
      if (num === _0n2)
        return acc;
      scratch[i4] = mod(acc * scratch[i4], p4);
      return mod(acc * num, p4);
    }, inverted);
    return scratch;
  }
  function bits2int_2(bytes2) {
    const delta = bytes2.length * 8 - groupLen * 8;
    const num = bytesToNumber(bytes2);
    return delta > 0 ? num >> BigInt(delta) : num;
  }
  function truncateHash(hash2, truncateOnly = false) {
    const h4 = bits2int_2(hash2);
    if (truncateOnly)
      return h4;
    const { n: n4 } = CURVE;
    return h4 >= n4 ? h4 - n4 : h4;
  }
  var _sha256Sync;
  var _hmacSha256Sync;
  var HmacDrbg = class {
    constructor(hashLen2, qByteLen) {
      this.hashLen = hashLen2;
      this.qByteLen = qByteLen;
      if (typeof hashLen2 !== "number" || hashLen2 < 2)
        throw new Error("hashLen must be a number");
      if (typeof qByteLen !== "number" || qByteLen < 2)
        throw new Error("qByteLen must be a number");
      this.v = new Uint8Array(hashLen2).fill(1);
      this.k = new Uint8Array(hashLen2).fill(0);
      this.counter = 0;
    }
    hmac(...values) {
      return utils.hmacSha256(this.k, ...values);
    }
    hmacSync(...values) {
      return _hmacSha256Sync(this.k, ...values);
    }
    checkSync() {
      if (typeof _hmacSha256Sync !== "function")
        throw new ShaError("hmacSha256Sync needs to be set");
    }
    incr() {
      if (this.counter >= 1e3)
        throw new Error("Tried 1,000 k values for sign(), all were invalid");
      this.counter += 1;
    }
    async reseed(seed = new Uint8Array()) {
      this.k = await this.hmac(this.v, Uint8Array.from([0]), seed);
      this.v = await this.hmac(this.v);
      if (seed.length === 0)
        return;
      this.k = await this.hmac(this.v, Uint8Array.from([1]), seed);
      this.v = await this.hmac(this.v);
    }
    reseedSync(seed = new Uint8Array()) {
      this.checkSync();
      this.k = this.hmacSync(this.v, Uint8Array.from([0]), seed);
      this.v = this.hmacSync(this.v);
      if (seed.length === 0)
        return;
      this.k = this.hmacSync(this.v, Uint8Array.from([1]), seed);
      this.v = this.hmacSync(this.v);
    }
    async generate() {
      this.incr();
      let len = 0;
      const out = [];
      while (len < this.qByteLen) {
        this.v = await this.hmac(this.v);
        const sl = this.v.slice();
        out.push(sl);
        len += this.v.length;
      }
      return concatBytes(...out);
    }
    generateSync() {
      this.checkSync();
      this.incr();
      let len = 0;
      const out = [];
      while (len < this.qByteLen) {
        this.v = this.hmacSync(this.v);
        const sl = this.v.slice();
        out.push(sl);
        len += this.v.length;
      }
      return concatBytes(...out);
    }
  };
  function isWithinCurveOrder(num) {
    return _0n2 < num && num < CURVE.n;
  }
  function isValidFieldElement(num) {
    return _0n2 < num && num < CURVE.P;
  }
  function kmdToSig(kBytes, m4, d4, lowS = true) {
    const { n: n4 } = CURVE;
    const k3 = truncateHash(kBytes, true);
    if (!isWithinCurveOrder(k3))
      return;
    const kinv = invert(k3, n4);
    const q2 = Point.BASE.multiply(k3);
    const r4 = mod(q2.x, n4);
    if (r4 === _0n2)
      return;
    const s4 = mod(kinv * mod(m4 + d4 * r4, n4), n4);
    if (s4 === _0n2)
      return;
    let sig = new Signature(r4, s4);
    let recovery = (q2.x === sig.r ? 0 : 2) | Number(q2.y & _1n2);
    if (lowS && sig.hasHighS()) {
      sig = sig.normalizeS();
      recovery ^= 1;
    }
    return { sig, recovery };
  }
  function normalizePrivateKey(key) {
    let num;
    if (typeof key === "bigint") {
      num = key;
    } else if (typeof key === "number" && Number.isSafeInteger(key) && key > 0) {
      num = BigInt(key);
    } else if (typeof key === "string") {
      if (key.length !== 2 * groupLen)
        throw new Error("Expected 32 bytes of private key");
      num = hexToNumber(key);
    } else if (key instanceof Uint8Array) {
      if (key.length !== groupLen)
        throw new Error("Expected 32 bytes of private key");
      num = bytesToNumber(key);
    } else {
      throw new TypeError("Expected valid private key");
    }
    if (!isWithinCurveOrder(num))
      throw new Error("Expected private key: 0 < key < n");
    return num;
  }
  function normalizePublicKey(publicKey) {
    if (publicKey instanceof Point) {
      publicKey.assertValidity();
      return publicKey;
    } else {
      return Point.fromHex(publicKey);
    }
  }
  function normalizeSignature(signature) {
    if (signature instanceof Signature) {
      signature.assertValidity();
      return signature;
    }
    try {
      return Signature.fromDER(signature);
    } catch (error) {
      return Signature.fromCompact(signature);
    }
  }
  function getPublicKey(privateKey, isCompressed = false) {
    return Point.fromPrivateKey(privateKey).toRawBytes(isCompressed);
  }
  function recoverPublicKey(msgHash, signature, recovery, isCompressed = false) {
    return Point.fromSignature(msgHash, signature, recovery).toRawBytes(isCompressed);
  }
  function isProbPub(item) {
    const arr = item instanceof Uint8Array;
    const str = typeof item === "string";
    const len = (arr || str) && item.length;
    if (arr)
      return len === compressedLen || len === uncompressedLen;
    if (str)
      return len === compressedLen * 2 || len === uncompressedLen * 2;
    if (item instanceof Point)
      return true;
    return false;
  }
  function getSharedSecret(privateA, publicB, isCompressed = false) {
    if (isProbPub(privateA))
      throw new TypeError("getSharedSecret: first arg must be private key");
    if (!isProbPub(publicB))
      throw new TypeError("getSharedSecret: second arg must be public key");
    const b5 = normalizePublicKey(publicB);
    b5.assertValidity();
    return b5.multiply(normalizePrivateKey(privateA)).toRawBytes(isCompressed);
  }
  function bits2int(bytes2) {
    const slice = bytes2.length > fieldLen ? bytes2.slice(0, fieldLen) : bytes2;
    return bytesToNumber(slice);
  }
  function bits2octets(bytes2) {
    const z1 = bits2int(bytes2);
    const z22 = mod(z1, CURVE.n);
    return int2octets(z22 < _0n2 ? z1 : z22);
  }
  function int2octets(num) {
    return numTo32b(num);
  }
  function initSigArgs(msgHash, privateKey, extraEntropy) {
    if (msgHash == null)
      throw new Error(`sign: expected valid message hash, not "${msgHash}"`);
    const h1 = ensureBytes(msgHash);
    const d4 = normalizePrivateKey(privateKey);
    const seedArgs = [int2octets(d4), bits2octets(h1)];
    if (extraEntropy != null) {
      if (extraEntropy === true)
        extraEntropy = utils.randomBytes(fieldLen);
      const e4 = ensureBytes(extraEntropy);
      if (e4.length !== fieldLen)
        throw new Error(`sign: Expected ${fieldLen} bytes of extra data`);
      seedArgs.push(e4);
    }
    const seed = concatBytes(...seedArgs);
    const m4 = bits2int(h1);
    return { seed, m: m4, d: d4 };
  }
  function finalizeSig(recSig, opts) {
    const { sig, recovery } = recSig;
    const { der, recovered } = Object.assign({ canonical: true, der: true }, opts);
    const hashed = der ? sig.toDERRawBytes() : sig.toCompactRawBytes();
    return recovered ? [hashed, recovery] : hashed;
  }
  function signSync(msgHash, privKey, opts = {}) {
    const { seed, m: m4, d: d4 } = initSigArgs(msgHash, privKey, opts.extraEntropy);
    const drbg = new HmacDrbg(hashLen, groupLen);
    drbg.reseedSync(seed);
    let sig;
    while (!(sig = kmdToSig(drbg.generateSync(), m4, d4, opts.canonical)))
      drbg.reseedSync();
    return finalizeSig(sig, opts);
  }
  Point.BASE._setWindowSize(8);
  var crypto3 = {
    node: nodeCrypto,
    web: typeof self === "object" && "crypto" in self ? self.crypto : void 0
  };
  var TAGGED_HASH_PREFIXES = {};
  var utils = {
    bytesToHex,
    hexToBytes,
    concatBytes,
    mod,
    invert,
    isValidPrivateKey(privateKey) {
      try {
        normalizePrivateKey(privateKey);
        return true;
      } catch (error) {
        return false;
      }
    },
    _bigintTo32Bytes: numTo32b,
    _normalizePrivateKey: normalizePrivateKey,
    hashToPrivateKey: (hash2) => {
      hash2 = ensureBytes(hash2);
      const minLen = groupLen + 8;
      if (hash2.length < minLen || hash2.length > 1024) {
        throw new Error(`Expected valid bytes of private key as per FIPS 186`);
      }
      const num = mod(bytesToNumber(hash2), CURVE.n - _1n2) + _1n2;
      return numTo32b(num);
    },
    randomBytes: (bytesLength = 32) => {
      if (crypto3.web) {
        return crypto3.web.getRandomValues(new Uint8Array(bytesLength));
      } else if (crypto3.node) {
        const { randomBytes } = crypto3.node;
        return Uint8Array.from(randomBytes(bytesLength));
      } else {
        throw new Error("The environment doesn't have randomBytes function");
      }
    },
    randomPrivateKey: () => utils.hashToPrivateKey(utils.randomBytes(groupLen + 8)),
    precompute(windowSize = 8, point = Point.BASE) {
      const cached = point === Point.BASE ? point : new Point(point.x, point.y);
      cached._setWindowSize(windowSize);
      cached.multiply(_3n);
      return cached;
    },
    sha256: async (...messages) => {
      if (crypto3.web) {
        const buffer = await crypto3.web.subtle.digest("SHA-256", concatBytes(...messages));
        return new Uint8Array(buffer);
      } else if (crypto3.node) {
        const { createHash } = crypto3.node;
        const hash2 = createHash("sha256");
        messages.forEach((m4) => hash2.update(m4));
        return Uint8Array.from(hash2.digest());
      } else {
        throw new Error("The environment doesn't have sha256 function");
      }
    },
    hmacSha256: async (key, ...messages) => {
      if (crypto3.web) {
        const ckey = await crypto3.web.subtle.importKey("raw", key, { name: "HMAC", hash: { name: "SHA-256" } }, false, ["sign"]);
        const message = concatBytes(...messages);
        const buffer = await crypto3.web.subtle.sign("HMAC", ckey, message);
        return new Uint8Array(buffer);
      } else if (crypto3.node) {
        const { createHmac: createHmac2 } = crypto3.node;
        const hash2 = createHmac2("sha256", key);
        messages.forEach((m4) => hash2.update(m4));
        return Uint8Array.from(hash2.digest());
      } else {
        throw new Error("The environment doesn't have hmac-sha256 function");
      }
    },
    sha256Sync: void 0,
    hmacSha256Sync: void 0,
    taggedHash: async (tag, ...messages) => {
      let tagP = TAGGED_HASH_PREFIXES[tag];
      if (tagP === void 0) {
        const tagH = await utils.sha256(Uint8Array.from(tag, (c4) => c4.charCodeAt(0)));
        tagP = concatBytes(tagH, tagH);
        TAGGED_HASH_PREFIXES[tag] = tagP;
      }
      return utils.sha256(tagP, ...messages);
    },
    taggedHashSync: (tag, ...messages) => {
      if (typeof _sha256Sync !== "function")
        throw new ShaError("sha256Sync is undefined, you need to set it");
      let tagP = TAGGED_HASH_PREFIXES[tag];
      if (tagP === void 0) {
        const tagH = _sha256Sync(Uint8Array.from(tag, (c4) => c4.charCodeAt(0)));
        tagP = concatBytes(tagH, tagH);
        TAGGED_HASH_PREFIXES[tag] = tagP;
      }
      return _sha256Sync(tagP, ...messages);
    },
    _JacobianPoint: JacobianPoint
  };
  Object.defineProperties(utils, {
    sha256Sync: {
      configurable: false,
      get() {
        return _sha256Sync;
      },
      set(val) {
        if (!_sha256Sync)
          _sha256Sync = val;
      }
    },
    hmacSha256Sync: {
      configurable: false,
      get() {
        return _hmacSha256Sync;
      },
      set(val) {
        if (!_hmacSha256Sync)
          _hmacSha256Sync = val;
      }
    }
  });

  // node_modules/ethers/lib.esm/constants/addresses.js
  var ZeroAddress = "0x0000000000000000000000000000000000000000";

  // node_modules/ethers/lib.esm/constants/hashes.js
  var ZeroHash = "0x0000000000000000000000000000000000000000000000000000000000000000";

  // node_modules/ethers/lib.esm/crypto/signature.js
  var BN_04 = BigInt(0);
  var BN_13 = BigInt(1);
  var BN_2 = BigInt(2);
  var BN_27 = BigInt(27);
  var BN_28 = BigInt(28);
  var BN_35 = BigInt(35);
  var _guard3 = {};
  function toUint256(value) {
    return zeroPadValue(toBeArray(value), 32);
  }
  var Signature2 = class _Signature {
    #r;
    #s;
    #v;
    #networkV;
    /**
     *  The ``r`` value for a signautre.
     *
     *  This represents the ``x`` coordinate of a "reference" or
     *  challenge point, from which the ``y`` can be computed.
     */
    get r() {
      return this.#r;
    }
    set r(value) {
      assertArgument(dataLength(value) === 32, "invalid r", "value", value);
      this.#r = hexlify(value);
    }
    /**
     *  The ``s`` value for a signature.
     */
    get s() {
      return this.#s;
    }
    set s(_value) {
      assertArgument(dataLength(_value) === 32, "invalid s", "value", _value);
      const value = hexlify(_value);
      assertArgument(parseInt(value.substring(0, 3)) < 8, "non-canonical s", "value", value);
      this.#s = value;
    }
    /**
     *  The ``v`` value for a signature.
     *
     *  Since a given ``x`` value for ``r`` has two possible values for
     *  its correspondin ``y``, the ``v`` indicates which of the two ``y``
     *  values to use.
     *
     *  It is normalized to the values ``27`` or ``28`` for legacy
     *  purposes.
     */
    get v() {
      return this.#v;
    }
    set v(value) {
      const v3 = getNumber(value, "value");
      assertArgument(v3 === 27 || v3 === 28, "invalid v", "v", value);
      this.#v = v3;
    }
    /**
     *  The EIP-155 ``v`` for legacy transactions. For non-legacy
     *  transactions, this value is ``null``.
     */
    get networkV() {
      return this.#networkV;
    }
    /**
     *  The chain ID for EIP-155 legacy transactions. For non-legacy
     *  transactions, this value is ``null``.
     */
    get legacyChainId() {
      const v3 = this.networkV;
      if (v3 == null) {
        return null;
      }
      return _Signature.getChainId(v3);
    }
    /**
     *  The ``yParity`` for the signature.
     *
     *  See ``v`` for more details on how this value is used.
     */
    get yParity() {
      return this.v === 27 ? 0 : 1;
    }
    /**
     *  The [[link-eip-2098]] compact representation of the ``yParity``
     *  and ``s`` compacted into a single ``bytes32``.
     */
    get yParityAndS() {
      const yParityAndS = getBytes(this.s);
      if (this.yParity) {
        yParityAndS[0] |= 128;
      }
      return hexlify(yParityAndS);
    }
    /**
     *  The [[link-eip-2098]] compact representation.
     */
    get compactSerialized() {
      return concat([this.r, this.yParityAndS]);
    }
    /**
     *  The serialized representation.
     */
    get serialized() {
      return concat([this.r, this.s, this.yParity ? "0x1c" : "0x1b"]);
    }
    /**
     *  @private
     */
    constructor(guard, r4, s4, v3) {
      assertPrivate(guard, _guard3, "Signature");
      this.#r = r4;
      this.#s = s4;
      this.#v = v3;
      this.#networkV = null;
    }
    [Symbol.for("nodejs.util.inspect.custom")]() {
      return `Signature { r: "${this.r}", s: "${this.s}", yParity: ${this.yParity}, networkV: ${this.networkV} }`;
    }
    /**
     *  Returns a new identical [[Signature]].
     */
    clone() {
      const clone = new _Signature(_guard3, this.r, this.s, this.v);
      if (this.networkV) {
        clone.#networkV = this.networkV;
      }
      return clone;
    }
    /**
     *  Returns a representation that is compatible with ``JSON.stringify``.
     */
    toJSON() {
      const networkV = this.networkV;
      return {
        _type: "signature",
        networkV: networkV != null ? networkV.toString() : null,
        r: this.r,
        s: this.s,
        v: this.v
      };
    }
    /**
     *  Compute the chain ID from the ``v`` in a legacy EIP-155 transactions.
     *
     *  @example:
     *    Signature.getChainId(45)
     *    //_result:
     *
     *    Signature.getChainId(46)
     *    //_result:
     */
    static getChainId(v3) {
      const bv = getBigInt(v3, "v");
      if (bv == BN_27 || bv == BN_28) {
        return BN_04;
      }
      assertArgument(bv >= BN_35, "invalid EIP-155 v", "v", v3);
      return (bv - BN_35) / BN_2;
    }
    /**
     *  Compute the ``v`` for a chain ID for a legacy EIP-155 transactions.
     *
     *  Legacy transactions which use [[link-eip-155]] hijack the ``v``
     *  property to include the chain ID.
     *
     *  @example:
     *    Signature.getChainIdV(5, 27)
     *    //_result:
     *
     *    Signature.getChainIdV(5, 28)
     *    //_result:
     *
     */
    static getChainIdV(chainId, v3) {
      return getBigInt(chainId) * BN_2 + BigInt(35 + v3 - 27);
    }
    /**
     *  Compute the normalized legacy transaction ``v`` from a ``yParirty``,
     *  a legacy transaction ``v`` or a legacy [[link-eip-155]] transaction.
     *
     *  @example:
     *    // The values 0 and 1 imply v is actually yParity
     *    Signature.getNormalizedV(0)
     *    //_result:
     *
     *    // Legacy non-EIP-1559 transaction (i.e. 27 or 28)
     *    Signature.getNormalizedV(27)
     *    //_result:
     *
     *    // Legacy EIP-155 transaction (i.e. >= 35)
     *    Signature.getNormalizedV(46)
     *    //_result:
     *
     *    // Invalid values throw
     *    Signature.getNormalizedV(5)
     *    //_error:
     */
    static getNormalizedV(v3) {
      const bv = getBigInt(v3);
      if (bv === BN_04 || bv === BN_27) {
        return 27;
      }
      if (bv === BN_13 || bv === BN_28) {
        return 28;
      }
      assertArgument(bv >= BN_35, "invalid v", "v", v3);
      return bv & BN_13 ? 27 : 28;
    }
    /**
     *  Creates a new [[Signature]].
     *
     *  If no %%sig%% is provided, a new [[Signature]] is created
     *  with default values.
     *
     *  If %%sig%% is a string, it is parsed.
     */
    static from(sig) {
      function assertError(check, message) {
        assertArgument(check, message, "signature", sig);
      }
      ;
      if (sig == null) {
        return new _Signature(_guard3, ZeroHash, ZeroHash, 27);
      }
      if (typeof sig === "string") {
        const bytes2 = getBytes(sig, "signature");
        if (bytes2.length === 64) {
          const r5 = hexlify(bytes2.slice(0, 32));
          const s5 = bytes2.slice(32, 64);
          const v4 = s5[0] & 128 ? 28 : 27;
          s5[0] &= 127;
          return new _Signature(_guard3, r5, hexlify(s5), v4);
        }
        if (bytes2.length === 65) {
          const r5 = hexlify(bytes2.slice(0, 32));
          const s5 = bytes2.slice(32, 64);
          assertError((s5[0] & 128) === 0, "non-canonical s");
          const v4 = _Signature.getNormalizedV(bytes2[64]);
          return new _Signature(_guard3, r5, hexlify(s5), v4);
        }
        assertError(false, "invalid raw signature length");
      }
      if (sig instanceof _Signature) {
        return sig.clone();
      }
      const _r = sig.r;
      assertError(_r != null, "missing r");
      const r4 = toUint256(_r);
      const s4 = function(s5, yParityAndS) {
        if (s5 != null) {
          return toUint256(s5);
        }
        if (yParityAndS != null) {
          assertError(isHexString(yParityAndS, 32), "invalid yParityAndS");
          const bytes2 = getBytes(yParityAndS);
          bytes2[0] &= 127;
          return hexlify(bytes2);
        }
        assertError(false, "missing s");
      }(sig.s, sig.yParityAndS);
      assertError((getBytes(s4)[0] & 128) == 0, "non-canonical s");
      const { networkV, v: v3 } = function(_v, yParityAndS, yParity) {
        if (_v != null) {
          const v4 = getBigInt(_v);
          return {
            networkV: v4 >= BN_35 ? v4 : void 0,
            v: _Signature.getNormalizedV(v4)
          };
        }
        if (yParityAndS != null) {
          assertError(isHexString(yParityAndS, 32), "invalid yParityAndS");
          return { v: getBytes(yParityAndS)[0] & 128 ? 28 : 27 };
        }
        if (yParity != null) {
          switch (getNumber(yParity, "sig.yParity")) {
            case 0:
              return { v: 27 };
            case 1:
              return { v: 28 };
          }
          assertError(false, "invalid yParity");
        }
        assertError(false, "missing v");
      }(sig.v, sig.yParityAndS, sig.yParity);
      const result = new _Signature(_guard3, r4, s4, v3);
      if (networkV) {
        result.#networkV = networkV;
      }
      assertError(sig.yParity == null || getNumber(sig.yParity, "sig.yParity") === result.yParity, "yParity mismatch");
      assertError(sig.yParityAndS == null || sig.yParityAndS === result.yParityAndS, "yParityAndS mismatch");
      return result;
    }
  };

  // node_modules/ethers/lib.esm/crypto/signing-key.js
  utils.hmacSha256Sync = function(key, ...messages) {
    return getBytes(computeHmac("sha256", key, concat(messages)));
  };
  var SigningKey = class _SigningKey {
    #privateKey;
    /**
     *  Creates a new **SigningKey** for %%privateKey%%.
     */
    constructor(privateKey) {
      assertArgument(dataLength(privateKey) === 32, "invalid private key", "privateKey", "[REDACTED]");
      this.#privateKey = hexlify(privateKey);
    }
    /**
     *  The private key.
     */
    get privateKey() {
      return this.#privateKey;
    }
    /**
     *  The uncompressed public key.
     *
     * This will always begin with the prefix ``0x04`` and be 132
     * characters long (the ``0x`` prefix and 130 hexadecimal nibbles).
     */
    get publicKey() {
      return _SigningKey.computePublicKey(this.#privateKey);
    }
    /**
     *  The compressed public key.
     *
     *  This will always begin with either the prefix ``0x02`` or ``0x03``
     *  and be 68 characters long (the ``0x`` prefix and 33 hexadecimal
     *  nibbles)
     */
    get compressedPublicKey() {
      return _SigningKey.computePublicKey(this.#privateKey, true);
    }
    /**
     *  Return the signature of the signed %%digest%%.
     */
    sign(digest) {
      assertArgument(dataLength(digest) === 32, "invalid digest length", "digest", digest);
      const [sigDer, recid] = signSync(getBytesCopy(digest), getBytesCopy(this.#privateKey), {
        recovered: true,
        canonical: true
      });
      const sig = Signature.fromHex(sigDer);
      return Signature2.from({
        r: toBeHex("0x" + sig.r.toString(16), 32),
        s: toBeHex("0x" + sig.s.toString(16), 32),
        v: recid ? 28 : 27
      });
    }
    /**
     *  Returns the [[link-wiki-ecdh]] shared secret between this
     *  private key and the %%other%% key.
     *
     *  The %%other%% key may be any type of key, a raw public key,
     *  a compressed/uncompressed pubic key or aprivate key.
     *
     *  Best practice is usually to use a cryptographic hash on the
     *  returned value before using it as a symetric secret.
     *
     *  @example:
     *    sign1 = new SigningKey(id("some-secret-1"))
     *    sign2 = new SigningKey(id("some-secret-2"))
     *
     *    // Notice that privA.computeSharedSecret(pubB)...
     *    sign1.computeSharedSecret(sign2.publicKey)
     *    //_result:
     *
     *    // ...is equal to privB.computeSharedSecret(pubA).
     *    sign2.computeSharedSecret(sign1.publicKey)
     *    //_result:
     */
    computeSharedSecret(other) {
      const pubKey = _SigningKey.computePublicKey(other);
      return hexlify(getSharedSecret(getBytesCopy(this.#privateKey), getBytes(pubKey)));
    }
    /**
     *  Compute the public key for %%key%%, optionally %%compressed%%.
     *
     *  The %%key%% may be any type of key, a raw public key, a
     *  compressed/uncompressed public key or private key.
     *
     *  @example:
     *    sign = new SigningKey(id("some-secret"));
     *
     *    // Compute the uncompressed public key for a private key
     *    SigningKey.computePublicKey(sign.privateKey)
     *    //_result:
     *
     *    // Compute the compressed public key for a private key
     *    SigningKey.computePublicKey(sign.privateKey, true)
     *    //_result:
     *
     *    // Compute the uncompressed public key
     *    SigningKey.computePublicKey(sign.publicKey, false);
     *    //_result:
     *
     *    // Compute the Compressed a public key
     *    SigningKey.computePublicKey(sign.publicKey, true);
     *    //_result:
     */
    static computePublicKey(key, compressed) {
      let bytes2 = getBytes(key, "key");
      if (bytes2.length === 32) {
        const pubKey = getPublicKey(bytes2, !!compressed);
        return hexlify(pubKey);
      }
      if (bytes2.length === 64) {
        const pub = new Uint8Array(65);
        pub[0] = 4;
        pub.set(bytes2, 1);
        bytes2 = pub;
      }
      const point = Point.fromHex(bytes2);
      return hexlify(point.toRawBytes(compressed));
    }
    /**
     *  Returns the public key for the private key which produced the
     *  %%signature%% for the given %%digest%%.
     *
     *  @example:
     *    key = new SigningKey(id("some-secret"))
     *    digest = id("hello world")
     *    sig = key.sign(digest)
     *
     *    // Notice the signer public key...
     *    key.publicKey
     *    //_result:
     *
     *    // ...is equal to the recovered public key
     *    SigningKey.recoverPublicKey(digest, sig)
     *    //_result:
     *
     */
    static recoverPublicKey(digest, signature) {
      assertArgument(dataLength(digest) === 32, "invalid digest length", "digest", digest);
      const sig = Signature2.from(signature);
      const der = Signature.fromCompact(getBytesCopy(concat([sig.r, sig.s]))).toDERRawBytes();
      const pubKey = recoverPublicKey(getBytesCopy(digest), der, sig.yParity);
      assertArgument(pubKey != null, "invalid signature for digest", "signature", signature);
      return hexlify(pubKey);
    }
    /**
     *  Returns the point resulting from adding the ellipic curve points
     *  %%p0%% and %%p1%%.
     *
     *  This is not a common function most developers should require, but
     *  can be useful for certain privacy-specific techniques.
     *
     *  For example, it is used by [[HDNodeWallet]] to compute child
     *  addresses from parent public keys and chain codes.
     */
    static addPoints(p0, p1, compressed) {
      const pub0 = Point.fromHex(_SigningKey.computePublicKey(p0).substring(2));
      const pub1 = Point.fromHex(_SigningKey.computePublicKey(p1).substring(2));
      return "0x" + pub0.add(pub1).toHex(!!compressed);
    }
  };

  // node_modules/ethers/lib.esm/address/address.js
  var BN_05 = BigInt(0);
  var BN_36 = BigInt(36);
  function getChecksumAddress(address) {
    address = address.toLowerCase();
    const chars = address.substring(2).split("");
    const expanded = new Uint8Array(40);
    for (let i4 = 0; i4 < 40; i4++) {
      expanded[i4] = chars[i4].charCodeAt(0);
    }
    const hashed = getBytes(keccak256(expanded));
    for (let i4 = 0; i4 < 40; i4 += 2) {
      if (hashed[i4 >> 1] >> 4 >= 8) {
        chars[i4] = chars[i4].toUpperCase();
      }
      if ((hashed[i4 >> 1] & 15) >= 8) {
        chars[i4 + 1] = chars[i4 + 1].toUpperCase();
      }
    }
    return "0x" + chars.join("");
  }
  var ibanLookup = {};
  for (let i4 = 0; i4 < 10; i4++) {
    ibanLookup[String(i4)] = String(i4);
  }
  for (let i4 = 0; i4 < 26; i4++) {
    ibanLookup[String.fromCharCode(65 + i4)] = String(10 + i4);
  }
  var safeDigits = 15;
  function ibanChecksum(address) {
    address = address.toUpperCase();
    address = address.substring(4) + address.substring(0, 2) + "00";
    let expanded = address.split("").map((c4) => {
      return ibanLookup[c4];
    }).join("");
    while (expanded.length >= safeDigits) {
      let block = expanded.substring(0, safeDigits);
      expanded = parseInt(block, 10) % 97 + expanded.substring(block.length);
    }
    let checksum = String(98 - parseInt(expanded, 10) % 97);
    while (checksum.length < 2) {
      checksum = "0" + checksum;
    }
    return checksum;
  }
  var Base36 = function() {
    ;
    const result = {};
    for (let i4 = 0; i4 < 36; i4++) {
      const key = "0123456789abcdefghijklmnopqrstuvwxyz"[i4];
      result[key] = BigInt(i4);
    }
    return result;
  }();
  function fromBase36(value) {
    value = value.toLowerCase();
    let result = BN_05;
    for (let i4 = 0; i4 < value.length; i4++) {
      result = result * BN_36 + Base36[value[i4]];
    }
    return result;
  }
  function getAddress(address) {
    assertArgument(typeof address === "string", "invalid address", "address", address);
    if (address.match(/^(0x)?[0-9a-fA-F]{40}$/)) {
      if (!address.startsWith("0x")) {
        address = "0x" + address;
      }
      const result = getChecksumAddress(address);
      assertArgument(!address.match(/([A-F].*[a-f])|([a-f].*[A-F])/) || result === address, "bad address checksum", "address", address);
      return result;
    }
    if (address.match(/^XE[0-9]{2}[0-9A-Za-z]{30,31}$/)) {
      assertArgument(address.substring(2, 4) === ibanChecksum(address), "bad icap checksum", "address", address);
      let result = fromBase36(address.substring(4)).toString(16);
      while (result.length < 40) {
        result = "0" + result;
      }
      return getChecksumAddress("0x" + result);
    }
    assertArgument(false, "invalid address", "address", address);
  }

  // node_modules/ethers/lib.esm/address/contract-address.js
  function getCreateAddress(tx) {
    const from = getAddress(tx.from);
    const nonce = getBigInt(tx.nonce, "tx.nonce");
    let nonceHex = nonce.toString(16);
    if (nonceHex === "0") {
      nonceHex = "0x";
    } else if (nonceHex.length % 2) {
      nonceHex = "0x0" + nonceHex;
    } else {
      nonceHex = "0x" + nonceHex;
    }
    return getAddress(dataSlice(keccak256(encodeRlp([from, nonceHex])), 12));
  }

  // node_modules/ethers/lib.esm/address/checks.js
  function isAddressable(value) {
    return value && typeof value.getAddress === "function";
  }
  async function checkAddress(target, promise) {
    const result = await promise;
    if (result == null || result === "0x0000000000000000000000000000000000000000") {
      assert(typeof target !== "string", "unconfigured name", "UNCONFIGURED_NAME", { value: target });
      assertArgument(false, "invalid AddressLike value; did not resolve to a value address", "target", target);
    }
    return getAddress(result);
  }
  function resolveAddress(target, resolver) {
    if (typeof target === "string") {
      if (target.match(/^0x[0-9a-f]{40}$/i)) {
        return getAddress(target);
      }
      assert(resolver != null, "ENS resolution requires a provider", "UNSUPPORTED_OPERATION", { operation: "resolveName" });
      return checkAddress(target, resolver.resolveName(target));
    } else if (isAddressable(target)) {
      return checkAddress(target, target.getAddress());
    } else if (target && typeof target.then === "function") {
      return checkAddress(target, target);
    }
    assertArgument(false, "unsupported addressable value", "target", target);
  }

  // node_modules/ethers/lib.esm/abi/typed.js
  var _gaurd = {};
  function n2(value, width) {
    let signed2 = false;
    if (width < 0) {
      signed2 = true;
      width *= -1;
    }
    return new Typed(_gaurd, `${signed2 ? "" : "u"}int${width}`, value, { signed: signed2, width });
  }
  function b3(value, size) {
    return new Typed(_gaurd, `bytes${size ? size : ""}`, value, { size });
  }
  var _typedSymbol = Symbol.for("_ethers_typed");
  var Typed = class _Typed {
    /**
     *  The type, as a Solidity-compatible type.
     */
    type;
    /**
     *  The actual value.
     */
    value;
    #options;
    /**
     *  @_ignore:
     */
    _typedSymbol;
    /**
     *  @_ignore:
     */
    constructor(gaurd, type, value, options) {
      if (options == null) {
        options = null;
      }
      assertPrivate(_gaurd, gaurd, "Typed");
      defineProperties(this, { _typedSymbol, type, value });
      this.#options = options;
      this.format();
    }
    /**
     *  Format the type as a Human-Readable type.
     */
    format() {
      if (this.type === "array") {
        throw new Error("");
      } else if (this.type === "dynamicArray") {
        throw new Error("");
      } else if (this.type === "tuple") {
        return `tuple(${this.value.map((v3) => v3.format()).join(",")})`;
      }
      return this.type;
    }
    /**
     *  The default value returned by this type.
     */
    defaultValue() {
      return 0;
    }
    /**
     *  The minimum value for numeric types.
     */
    minValue() {
      return 0;
    }
    /**
     *  The maximum value for numeric types.
     */
    maxValue() {
      return 0;
    }
    /**
     *  Returns ``true`` and provides a type guard is this is a [[TypedBigInt]].
     */
    isBigInt() {
      return !!this.type.match(/^u?int[0-9]+$/);
    }
    /**
     *  Returns ``true`` and provides a type guard is this is a [[TypedData]].
     */
    isData() {
      return this.type.startsWith("bytes");
    }
    /**
     *  Returns ``true`` and provides a type guard is this is a [[TypedString]].
     */
    isString() {
      return this.type === "string";
    }
    /**
     *  Returns the tuple name, if this is a tuple. Throws otherwise.
     */
    get tupleName() {
      if (this.type !== "tuple") {
        throw TypeError("not a tuple");
      }
      return this.#options;
    }
    // Returns the length of this type as an array
    // - `null` indicates the length is unforced, it could be dynamic
    // - `-1` indicates the length is dynamic
    // - any other value indicates it is a static array and is its length
    /**
     *  Returns the length of the array type or ``-1`` if it is dynamic.
     *
     *  Throws if the type is not an array.
     */
    get arrayLength() {
      if (this.type !== "array") {
        throw TypeError("not an array");
      }
      if (this.#options === true) {
        return -1;
      }
      if (this.#options === false) {
        return this.value.length;
      }
      return null;
    }
    /**
     *  Returns a new **Typed** of %%type%% with the %%value%%.
     */
    static from(type, value) {
      return new _Typed(_gaurd, type, value);
    }
    /**
     *  Return a new ``uint8`` type for %%v%%.
     */
    static uint8(v3) {
      return n2(v3, 8);
    }
    /**
     *  Return a new ``uint16`` type for %%v%%.
     */
    static uint16(v3) {
      return n2(v3, 16);
    }
    /**
     *  Return a new ``uint24`` type for %%v%%.
     */
    static uint24(v3) {
      return n2(v3, 24);
    }
    /**
     *  Return a new ``uint32`` type for %%v%%.
     */
    static uint32(v3) {
      return n2(v3, 32);
    }
    /**
     *  Return a new ``uint40`` type for %%v%%.
     */
    static uint40(v3) {
      return n2(v3, 40);
    }
    /**
     *  Return a new ``uint48`` type for %%v%%.
     */
    static uint48(v3) {
      return n2(v3, 48);
    }
    /**
     *  Return a new ``uint56`` type for %%v%%.
     */
    static uint56(v3) {
      return n2(v3, 56);
    }
    /**
     *  Return a new ``uint64`` type for %%v%%.
     */
    static uint64(v3) {
      return n2(v3, 64);
    }
    /**
     *  Return a new ``uint72`` type for %%v%%.
     */
    static uint72(v3) {
      return n2(v3, 72);
    }
    /**
     *  Return a new ``uint80`` type for %%v%%.
     */
    static uint80(v3) {
      return n2(v3, 80);
    }
    /**
     *  Return a new ``uint88`` type for %%v%%.
     */
    static uint88(v3) {
      return n2(v3, 88);
    }
    /**
     *  Return a new ``uint96`` type for %%v%%.
     */
    static uint96(v3) {
      return n2(v3, 96);
    }
    /**
     *  Return a new ``uint104`` type for %%v%%.
     */
    static uint104(v3) {
      return n2(v3, 104);
    }
    /**
     *  Return a new ``uint112`` type for %%v%%.
     */
    static uint112(v3) {
      return n2(v3, 112);
    }
    /**
     *  Return a new ``uint120`` type for %%v%%.
     */
    static uint120(v3) {
      return n2(v3, 120);
    }
    /**
     *  Return a new ``uint128`` type for %%v%%.
     */
    static uint128(v3) {
      return n2(v3, 128);
    }
    /**
     *  Return a new ``uint136`` type for %%v%%.
     */
    static uint136(v3) {
      return n2(v3, 136);
    }
    /**
     *  Return a new ``uint144`` type for %%v%%.
     */
    static uint144(v3) {
      return n2(v3, 144);
    }
    /**
     *  Return a new ``uint152`` type for %%v%%.
     */
    static uint152(v3) {
      return n2(v3, 152);
    }
    /**
     *  Return a new ``uint160`` type for %%v%%.
     */
    static uint160(v3) {
      return n2(v3, 160);
    }
    /**
     *  Return a new ``uint168`` type for %%v%%.
     */
    static uint168(v3) {
      return n2(v3, 168);
    }
    /**
     *  Return a new ``uint176`` type for %%v%%.
     */
    static uint176(v3) {
      return n2(v3, 176);
    }
    /**
     *  Return a new ``uint184`` type for %%v%%.
     */
    static uint184(v3) {
      return n2(v3, 184);
    }
    /**
     *  Return a new ``uint192`` type for %%v%%.
     */
    static uint192(v3) {
      return n2(v3, 192);
    }
    /**
     *  Return a new ``uint200`` type for %%v%%.
     */
    static uint200(v3) {
      return n2(v3, 200);
    }
    /**
     *  Return a new ``uint208`` type for %%v%%.
     */
    static uint208(v3) {
      return n2(v3, 208);
    }
    /**
     *  Return a new ``uint216`` type for %%v%%.
     */
    static uint216(v3) {
      return n2(v3, 216);
    }
    /**
     *  Return a new ``uint224`` type for %%v%%.
     */
    static uint224(v3) {
      return n2(v3, 224);
    }
    /**
     *  Return a new ``uint232`` type for %%v%%.
     */
    static uint232(v3) {
      return n2(v3, 232);
    }
    /**
     *  Return a new ``uint240`` type for %%v%%.
     */
    static uint240(v3) {
      return n2(v3, 240);
    }
    /**
     *  Return a new ``uint248`` type for %%v%%.
     */
    static uint248(v3) {
      return n2(v3, 248);
    }
    /**
     *  Return a new ``uint256`` type for %%v%%.
     */
    static uint256(v3) {
      return n2(v3, 256);
    }
    /**
     *  Return a new ``uint256`` type for %%v%%.
     */
    static uint(v3) {
      return n2(v3, 256);
    }
    /**
     *  Return a new ``int8`` type for %%v%%.
     */
    static int8(v3) {
      return n2(v3, -8);
    }
    /**
     *  Return a new ``int16`` type for %%v%%.
     */
    static int16(v3) {
      return n2(v3, -16);
    }
    /**
     *  Return a new ``int24`` type for %%v%%.
     */
    static int24(v3) {
      return n2(v3, -24);
    }
    /**
     *  Return a new ``int32`` type for %%v%%.
     */
    static int32(v3) {
      return n2(v3, -32);
    }
    /**
     *  Return a new ``int40`` type for %%v%%.
     */
    static int40(v3) {
      return n2(v3, -40);
    }
    /**
     *  Return a new ``int48`` type for %%v%%.
     */
    static int48(v3) {
      return n2(v3, -48);
    }
    /**
     *  Return a new ``int56`` type for %%v%%.
     */
    static int56(v3) {
      return n2(v3, -56);
    }
    /**
     *  Return a new ``int64`` type for %%v%%.
     */
    static int64(v3) {
      return n2(v3, -64);
    }
    /**
     *  Return a new ``int72`` type for %%v%%.
     */
    static int72(v3) {
      return n2(v3, -72);
    }
    /**
     *  Return a new ``int80`` type for %%v%%.
     */
    static int80(v3) {
      return n2(v3, -80);
    }
    /**
     *  Return a new ``int88`` type for %%v%%.
     */
    static int88(v3) {
      return n2(v3, -88);
    }
    /**
     *  Return a new ``int96`` type for %%v%%.
     */
    static int96(v3) {
      return n2(v3, -96);
    }
    /**
     *  Return a new ``int104`` type for %%v%%.
     */
    static int104(v3) {
      return n2(v3, -104);
    }
    /**
     *  Return a new ``int112`` type for %%v%%.
     */
    static int112(v3) {
      return n2(v3, -112);
    }
    /**
     *  Return a new ``int120`` type for %%v%%.
     */
    static int120(v3) {
      return n2(v3, -120);
    }
    /**
     *  Return a new ``int128`` type for %%v%%.
     */
    static int128(v3) {
      return n2(v3, -128);
    }
    /**
     *  Return a new ``int136`` type for %%v%%.
     */
    static int136(v3) {
      return n2(v3, -136);
    }
    /**
     *  Return a new ``int144`` type for %%v%%.
     */
    static int144(v3) {
      return n2(v3, -144);
    }
    /**
     *  Return a new ``int52`` type for %%v%%.
     */
    static int152(v3) {
      return n2(v3, -152);
    }
    /**
     *  Return a new ``int160`` type for %%v%%.
     */
    static int160(v3) {
      return n2(v3, -160);
    }
    /**
     *  Return a new ``int168`` type for %%v%%.
     */
    static int168(v3) {
      return n2(v3, -168);
    }
    /**
     *  Return a new ``int176`` type for %%v%%.
     */
    static int176(v3) {
      return n2(v3, -176);
    }
    /**
     *  Return a new ``int184`` type for %%v%%.
     */
    static int184(v3) {
      return n2(v3, -184);
    }
    /**
     *  Return a new ``int92`` type for %%v%%.
     */
    static int192(v3) {
      return n2(v3, -192);
    }
    /**
     *  Return a new ``int200`` type for %%v%%.
     */
    static int200(v3) {
      return n2(v3, -200);
    }
    /**
     *  Return a new ``int208`` type for %%v%%.
     */
    static int208(v3) {
      return n2(v3, -208);
    }
    /**
     *  Return a new ``int216`` type for %%v%%.
     */
    static int216(v3) {
      return n2(v3, -216);
    }
    /**
     *  Return a new ``int224`` type for %%v%%.
     */
    static int224(v3) {
      return n2(v3, -224);
    }
    /**
     *  Return a new ``int232`` type for %%v%%.
     */
    static int232(v3) {
      return n2(v3, -232);
    }
    /**
     *  Return a new ``int240`` type for %%v%%.
     */
    static int240(v3) {
      return n2(v3, -240);
    }
    /**
     *  Return a new ``int248`` type for %%v%%.
     */
    static int248(v3) {
      return n2(v3, -248);
    }
    /**
     *  Return a new ``int256`` type for %%v%%.
     */
    static int256(v3) {
      return n2(v3, -256);
    }
    /**
     *  Return a new ``int256`` type for %%v%%.
     */
    static int(v3) {
      return n2(v3, -256);
    }
    /**
     *  Return a new ``bytes1`` type for %%v%%.
     */
    static bytes1(v3) {
      return b3(v3, 1);
    }
    /**
     *  Return a new ``bytes2`` type for %%v%%.
     */
    static bytes2(v3) {
      return b3(v3, 2);
    }
    /**
     *  Return a new ``bytes3`` type for %%v%%.
     */
    static bytes3(v3) {
      return b3(v3, 3);
    }
    /**
     *  Return a new ``bytes4`` type for %%v%%.
     */
    static bytes4(v3) {
      return b3(v3, 4);
    }
    /**
     *  Return a new ``bytes5`` type for %%v%%.
     */
    static bytes5(v3) {
      return b3(v3, 5);
    }
    /**
     *  Return a new ``bytes6`` type for %%v%%.
     */
    static bytes6(v3) {
      return b3(v3, 6);
    }
    /**
     *  Return a new ``bytes7`` type for %%v%%.
     */
    static bytes7(v3) {
      return b3(v3, 7);
    }
    /**
     *  Return a new ``bytes8`` type for %%v%%.
     */
    static bytes8(v3) {
      return b3(v3, 8);
    }
    /**
     *  Return a new ``bytes9`` type for %%v%%.
     */
    static bytes9(v3) {
      return b3(v3, 9);
    }
    /**
     *  Return a new ``bytes10`` type for %%v%%.
     */
    static bytes10(v3) {
      return b3(v3, 10);
    }
    /**
     *  Return a new ``bytes11`` type for %%v%%.
     */
    static bytes11(v3) {
      return b3(v3, 11);
    }
    /**
     *  Return a new ``bytes12`` type for %%v%%.
     */
    static bytes12(v3) {
      return b3(v3, 12);
    }
    /**
     *  Return a new ``bytes13`` type for %%v%%.
     */
    static bytes13(v3) {
      return b3(v3, 13);
    }
    /**
     *  Return a new ``bytes14`` type for %%v%%.
     */
    static bytes14(v3) {
      return b3(v3, 14);
    }
    /**
     *  Return a new ``bytes15`` type for %%v%%.
     */
    static bytes15(v3) {
      return b3(v3, 15);
    }
    /**
     *  Return a new ``bytes16`` type for %%v%%.
     */
    static bytes16(v3) {
      return b3(v3, 16);
    }
    /**
     *  Return a new ``bytes17`` type for %%v%%.
     */
    static bytes17(v3) {
      return b3(v3, 17);
    }
    /**
     *  Return a new ``bytes18`` type for %%v%%.
     */
    static bytes18(v3) {
      return b3(v3, 18);
    }
    /**
     *  Return a new ``bytes19`` type for %%v%%.
     */
    static bytes19(v3) {
      return b3(v3, 19);
    }
    /**
     *  Return a new ``bytes20`` type for %%v%%.
     */
    static bytes20(v3) {
      return b3(v3, 20);
    }
    /**
     *  Return a new ``bytes21`` type for %%v%%.
     */
    static bytes21(v3) {
      return b3(v3, 21);
    }
    /**
     *  Return a new ``bytes22`` type for %%v%%.
     */
    static bytes22(v3) {
      return b3(v3, 22);
    }
    /**
     *  Return a new ``bytes23`` type for %%v%%.
     */
    static bytes23(v3) {
      return b3(v3, 23);
    }
    /**
     *  Return a new ``bytes24`` type for %%v%%.
     */
    static bytes24(v3) {
      return b3(v3, 24);
    }
    /**
     *  Return a new ``bytes25`` type for %%v%%.
     */
    static bytes25(v3) {
      return b3(v3, 25);
    }
    /**
     *  Return a new ``bytes26`` type for %%v%%.
     */
    static bytes26(v3) {
      return b3(v3, 26);
    }
    /**
     *  Return a new ``bytes27`` type for %%v%%.
     */
    static bytes27(v3) {
      return b3(v3, 27);
    }
    /**
     *  Return a new ``bytes28`` type for %%v%%.
     */
    static bytes28(v3) {
      return b3(v3, 28);
    }
    /**
     *  Return a new ``bytes29`` type for %%v%%.
     */
    static bytes29(v3) {
      return b3(v3, 29);
    }
    /**
     *  Return a new ``bytes30`` type for %%v%%.
     */
    static bytes30(v3) {
      return b3(v3, 30);
    }
    /**
     *  Return a new ``bytes31`` type for %%v%%.
     */
    static bytes31(v3) {
      return b3(v3, 31);
    }
    /**
     *  Return a new ``bytes32`` type for %%v%%.
     */
    static bytes32(v3) {
      return b3(v3, 32);
    }
    /**
     *  Return a new ``address`` type for %%v%%.
     */
    static address(v3) {
      return new _Typed(_gaurd, "address", v3);
    }
    /**
     *  Return a new ``bool`` type for %%v%%.
     */
    static bool(v3) {
      return new _Typed(_gaurd, "bool", !!v3);
    }
    /**
     *  Return a new ``bytes`` type for %%v%%.
     */
    static bytes(v3) {
      return new _Typed(_gaurd, "bytes", v3);
    }
    /**
     *  Return a new ``string`` type for %%v%%.
     */
    static string(v3) {
      return new _Typed(_gaurd, "string", v3);
    }
    /**
     *  Return a new ``array`` type for %%v%%, allowing %%dynamic%% length.
     */
    static array(v3, dynamic) {
      throw new Error("not implemented yet");
      return new _Typed(_gaurd, "array", v3, dynamic);
    }
    /**
     *  Return a new ``tuple`` type for %%v%%, with the optional %%name%%.
     */
    static tuple(v3, name) {
      throw new Error("not implemented yet");
      return new _Typed(_gaurd, "tuple", v3, name);
    }
    /**
     *  Return a new ``uint8`` type for %%v%%.
     */
    static overrides(v3) {
      return new _Typed(_gaurd, "overrides", Object.assign({}, v3));
    }
    /**
     *  Returns true only if %%value%% is a [[Typed]] instance.
     */
    static isTyped(value) {
      return value && typeof value === "object" && "_typedSymbol" in value && value._typedSymbol === _typedSymbol;
    }
    /**
     *  If the value is a [[Typed]] instance, validates the underlying value
     *  and returns it, otherwise returns value directly.
     *
     *  This is useful for functions that with to accept either a [[Typed]]
     *  object or values.
     */
    static dereference(value, type) {
      if (_Typed.isTyped(value)) {
        if (value.type !== type) {
          throw new Error(`invalid type: expecetd ${type}, got ${value.type}`);
        }
        return value.value;
      }
      return value;
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/address.js
  var AddressCoder = class extends Coder {
    constructor(localName) {
      super("address", "address", localName, false);
    }
    defaultValue() {
      return "0x0000000000000000000000000000000000000000";
    }
    encode(writer, _value) {
      let value = Typed.dereference(_value, "string");
      try {
        value = getAddress(value);
      } catch (error) {
        return this._throwError(error.message, _value);
      }
      return writer.writeValue(value);
    }
    decode(reader) {
      return getAddress(toBeHex(reader.readValue(), 20));
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/anonymous.js
  var AnonymousCoder = class extends Coder {
    coder;
    constructor(coder) {
      super(coder.name, coder.type, "_", coder.dynamic);
      this.coder = coder;
    }
    defaultValue() {
      return this.coder.defaultValue();
    }
    encode(writer, value) {
      return this.coder.encode(writer, value);
    }
    decode(reader) {
      return this.coder.decode(reader);
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/array.js
  function pack(writer, coders, values) {
    let arrayValues = [];
    if (Array.isArray(values)) {
      arrayValues = values;
    } else if (values && typeof values === "object") {
      let unique = {};
      arrayValues = coders.map((coder) => {
        const name = coder.localName;
        assert(name, "cannot encode object for signature with missing names", "INVALID_ARGUMENT", { argument: "values", info: { coder }, value: values });
        assert(!unique[name], "cannot encode object for signature with duplicate names", "INVALID_ARGUMENT", { argument: "values", info: { coder }, value: values });
        unique[name] = true;
        return values[name];
      });
    } else {
      assertArgument(false, "invalid tuple value", "tuple", values);
    }
    assertArgument(coders.length === arrayValues.length, "types/value length mismatch", "tuple", values);
    let staticWriter = new Writer();
    let dynamicWriter = new Writer();
    let updateFuncs = [];
    coders.forEach((coder, index) => {
      let value = arrayValues[index];
      if (coder.dynamic) {
        let dynamicOffset = dynamicWriter.length;
        coder.encode(dynamicWriter, value);
        let updateFunc = staticWriter.writeUpdatableValue();
        updateFuncs.push((baseOffset) => {
          updateFunc(baseOffset + dynamicOffset);
        });
      } else {
        coder.encode(staticWriter, value);
      }
    });
    updateFuncs.forEach((func) => {
      func(staticWriter.length);
    });
    let length = writer.appendWriter(staticWriter);
    length += writer.appendWriter(dynamicWriter);
    return length;
  }
  function unpack(reader, coders) {
    let values = [];
    let keys = [];
    let baseReader = reader.subReader(0);
    coders.forEach((coder) => {
      let value = null;
      if (coder.dynamic) {
        let offset = reader.readIndex();
        let offsetReader = baseReader.subReader(offset);
        try {
          value = coder.decode(offsetReader);
        } catch (error) {
          if (isError(error, "BUFFER_OVERRUN")) {
            throw error;
          }
          value = error;
          value.baseType = coder.name;
          value.name = coder.localName;
          value.type = coder.type;
        }
      } else {
        try {
          value = coder.decode(reader);
        } catch (error) {
          if (isError(error, "BUFFER_OVERRUN")) {
            throw error;
          }
          value = error;
          value.baseType = coder.name;
          value.name = coder.localName;
          value.type = coder.type;
        }
      }
      if (value == void 0) {
        throw new Error("investigate");
      }
      values.push(value);
      keys.push(coder.localName || null);
    });
    return Result.fromItems(values, keys);
  }
  var ArrayCoder = class extends Coder {
    coder;
    length;
    constructor(coder, length, localName) {
      const type = coder.type + "[" + (length >= 0 ? length : "") + "]";
      const dynamic = length === -1 || coder.dynamic;
      super("array", type, localName, dynamic);
      defineProperties(this, { coder, length });
    }
    defaultValue() {
      const defaultChild = this.coder.defaultValue();
      const result = [];
      for (let i4 = 0; i4 < this.length; i4++) {
        result.push(defaultChild);
      }
      return result;
    }
    encode(writer, _value) {
      const value = Typed.dereference(_value, "array");
      if (!Array.isArray(value)) {
        this._throwError("expected array value", value);
      }
      let count = this.length;
      if (count === -1) {
        count = value.length;
        writer.writeValue(value.length);
      }
      assertArgumentCount(value.length, count, "coder array" + (this.localName ? " " + this.localName : ""));
      let coders = [];
      for (let i4 = 0; i4 < value.length; i4++) {
        coders.push(this.coder);
      }
      return pack(writer, coders, value);
    }
    decode(reader) {
      let count = this.length;
      if (count === -1) {
        count = reader.readIndex();
        assert(count * WordSize <= reader.dataLength, "insufficient data length", "BUFFER_OVERRUN", { buffer: reader.bytes, offset: count * WordSize, length: reader.dataLength });
      }
      let coders = [];
      for (let i4 = 0; i4 < count; i4++) {
        coders.push(new AnonymousCoder(this.coder));
      }
      return unpack(reader, coders);
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/boolean.js
  var BooleanCoder = class extends Coder {
    constructor(localName) {
      super("bool", "bool", localName, false);
    }
    defaultValue() {
      return false;
    }
    encode(writer, _value) {
      const value = Typed.dereference(_value, "bool");
      return writer.writeValue(value ? 1 : 0);
    }
    decode(reader) {
      return !!reader.readValue();
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/bytes.js
  var DynamicBytesCoder = class extends Coder {
    constructor(type, localName) {
      super(type, type, localName, true);
    }
    defaultValue() {
      return "0x";
    }
    encode(writer, value) {
      value = getBytesCopy(value);
      let length = writer.writeValue(value.length);
      length += writer.writeBytes(value);
      return length;
    }
    decode(reader) {
      return reader.readBytes(reader.readIndex(), true);
    }
  };
  var BytesCoder = class extends DynamicBytesCoder {
    constructor(localName) {
      super("bytes", localName);
    }
    decode(reader) {
      return hexlify(super.decode(reader));
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/fixed-bytes.js
  var FixedBytesCoder = class extends Coder {
    size;
    constructor(size, localName) {
      let name = "bytes" + String(size);
      super(name, name, localName, false);
      defineProperties(this, { size }, { size: "number" });
    }
    defaultValue() {
      return "0x0000000000000000000000000000000000000000000000000000000000000000".substring(0, 2 + this.size * 2);
    }
    encode(writer, _value) {
      let data = getBytesCopy(Typed.dereference(_value, this.type));
      if (data.length !== this.size) {
        this._throwError("incorrect data length", _value);
      }
      return writer.writeBytes(data);
    }
    decode(reader) {
      return hexlify(reader.readBytes(this.size));
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/null.js
  var Empty = new Uint8Array([]);
  var NullCoder = class extends Coder {
    constructor(localName) {
      super("null", "", localName, false);
    }
    defaultValue() {
      return null;
    }
    encode(writer, value) {
      if (value != null) {
        this._throwError("not null", value);
      }
      return writer.writeBytes(Empty);
    }
    decode(reader) {
      reader.readBytes(0);
      return null;
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/number.js
  var BN_06 = BigInt(0);
  var BN_14 = BigInt(1);
  var BN_MAX_UINT256 = BigInt("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
  var NumberCoder = class extends Coder {
    size;
    signed;
    constructor(size, signed2, localName) {
      const name = (signed2 ? "int" : "uint") + size * 8;
      super(name, name, localName, false);
      defineProperties(this, { size, signed: signed2 }, { size: "number", signed: "boolean" });
    }
    defaultValue() {
      return 0;
    }
    encode(writer, _value) {
      let value = getBigInt(Typed.dereference(_value, this.type));
      let maxUintValue = mask(BN_MAX_UINT256, WordSize * 8);
      if (this.signed) {
        let bounds = mask(maxUintValue, this.size * 8 - 1);
        if (value > bounds || value < -(bounds + BN_14)) {
          this._throwError("value out-of-bounds", _value);
        }
        value = toTwos(value, 8 * WordSize);
      } else if (value < BN_06 || value > mask(maxUintValue, this.size * 8)) {
        this._throwError("value out-of-bounds", _value);
      }
      return writer.writeValue(value);
    }
    decode(reader) {
      let value = mask(reader.readValue(), this.size * 8);
      if (this.signed) {
        value = fromTwos(value, this.size * 8);
      }
      return value;
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/string.js
  var StringCoder = class extends DynamicBytesCoder {
    constructor(localName) {
      super("string", localName);
    }
    defaultValue() {
      return "";
    }
    encode(writer, _value) {
      return super.encode(writer, toUtf8Bytes(Typed.dereference(_value, "string")));
    }
    decode(reader) {
      return toUtf8String(super.decode(reader));
    }
  };

  // node_modules/ethers/lib.esm/abi/coders/tuple.js
  var TupleCoder = class extends Coder {
    coders;
    constructor(coders, localName) {
      let dynamic = false;
      const types = [];
      coders.forEach((coder) => {
        if (coder.dynamic) {
          dynamic = true;
        }
        types.push(coder.type);
      });
      const type = "tuple(" + types.join(",") + ")";
      super("tuple", type, localName, dynamic);
      defineProperties(this, { coders: Object.freeze(coders.slice()) });
    }
    defaultValue() {
      const values = [];
      this.coders.forEach((coder) => {
        values.push(coder.defaultValue());
      });
      const uniqueNames = this.coders.reduce((accum, coder) => {
        const name = coder.localName;
        if (name) {
          if (!accum[name]) {
            accum[name] = 0;
          }
          accum[name]++;
        }
        return accum;
      }, {});
      this.coders.forEach((coder, index) => {
        let name = coder.localName;
        if (!name || uniqueNames[name] !== 1) {
          return;
        }
        if (name === "length") {
          name = "_length";
        }
        if (values[name] != null) {
          return;
        }
        values[name] = values[index];
      });
      return Object.freeze(values);
    }
    encode(writer, _value) {
      const value = Typed.dereference(_value, "tuple");
      return pack(writer, this.coders, value);
    }
    decode(reader) {
      return unpack(reader, this.coders);
    }
  };

  // node_modules/ethers/lib.esm/hash/id.js
  function id(value) {
    return keccak256(toUtf8Bytes(value));
  }

  // node_modules/@adraffy/ens-normalize/dist/index.mjs
  function decode_arithmetic(bytes2) {
    let pos = 0;
    function u16() {
      return bytes2[pos++] << 8 | bytes2[pos++];
    }
    let symbol_count = u16();
    let total = 1;
    let acc = [0, 1];
    for (let i4 = 1; i4 < symbol_count; i4++) {
      acc.push(total += u16());
    }
    let skip = u16();
    let pos_payload = pos;
    pos += skip;
    let read_width = 0;
    let read_buffer = 0;
    function read_bit() {
      if (read_width == 0) {
        read_buffer = read_buffer << 8 | bytes2[pos++];
        read_width = 8;
      }
      return read_buffer >> --read_width & 1;
    }
    const N2 = 31;
    const FULL = 2 ** N2;
    const HALF = FULL >>> 1;
    const QRTR = HALF >> 1;
    const MASK = FULL - 1;
    let register = 0;
    for (let i4 = 0; i4 < N2; i4++)
      register = register << 1 | read_bit();
    let symbols = [];
    let low = 0;
    let range = FULL;
    while (true) {
      let value = Math.floor(((register - low + 1) * total - 1) / range);
      let start = 0;
      let end = symbol_count;
      while (end - start > 1) {
        let mid = start + end >>> 1;
        if (value < acc[mid]) {
          end = mid;
        } else {
          start = mid;
        }
      }
      if (start == 0)
        break;
      symbols.push(start);
      let a4 = low + Math.floor(range * acc[start] / total);
      let b5 = low + Math.floor(range * acc[start + 1] / total) - 1;
      while (((a4 ^ b5) & HALF) == 0) {
        register = register << 1 & MASK | read_bit();
        a4 = a4 << 1 & MASK;
        b5 = b5 << 1 & MASK | 1;
      }
      while (a4 & ~b5 & QRTR) {
        register = register & HALF | register << 1 & MASK >>> 1 | read_bit();
        a4 = a4 << 1 ^ HALF;
        b5 = (b5 ^ HALF) << 1 | HALF | 1;
      }
      low = a4;
      range = 1 + b5 - a4;
    }
    let offset = symbol_count - 4;
    return symbols.map((x2) => {
      switch (x2 - offset) {
        case 3:
          return offset + 65792 + (bytes2[pos_payload++] << 16 | bytes2[pos_payload++] << 8 | bytes2[pos_payload++]);
        case 2:
          return offset + 256 + (bytes2[pos_payload++] << 8 | bytes2[pos_payload++]);
        case 1:
          return offset + bytes2[pos_payload++];
        default:
          return x2 - 1;
      }
    });
  }
  function read_payload(v3) {
    let pos = 0;
    return () => v3[pos++];
  }
  function read_compressed_payload(s4) {
    return read_payload(decode_arithmetic(unsafe_atob(s4)));
  }
  function unsafe_atob(s4) {
    let lookup = [];
    [..."ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"].forEach((c4, i4) => lookup[c4.charCodeAt(0)] = i4);
    let n4 = s4.length;
    let ret = new Uint8Array(6 * n4 >> 3);
    for (let i4 = 0, pos = 0, width = 0, carry = 0; i4 < n4; i4++) {
      carry = carry << 6 | lookup[s4.charCodeAt(i4)];
      width += 6;
      if (width >= 8) {
        ret[pos++] = carry >> (width -= 8);
      }
    }
    return ret;
  }
  function signed(i4) {
    return i4 & 1 ? ~i4 >> 1 : i4 >> 1;
  }
  function read_deltas(n4, next) {
    let v3 = Array(n4);
    for (let i4 = 0, x2 = 0; i4 < n4; i4++)
      v3[i4] = x2 += signed(next());
    return v3;
  }
  function read_sorted(next, prev = 0) {
    let ret = [];
    while (true) {
      let x2 = next();
      let n4 = next();
      if (!n4)
        break;
      prev += x2;
      for (let i4 = 0; i4 < n4; i4++) {
        ret.push(prev + i4);
      }
      prev += n4 + 1;
    }
    return ret;
  }
  function read_sorted_arrays(next) {
    return read_array_while(() => {
      let v3 = read_sorted(next);
      if (v3.length)
        return v3;
    });
  }
  function read_mapped(next) {
    let ret = [];
    while (true) {
      let w3 = next();
      if (w3 == 0)
        break;
      ret.push(read_linear_table(w3, next));
    }
    while (true) {
      let w3 = next() - 1;
      if (w3 < 0)
        break;
      ret.push(read_replacement_table(w3, next));
    }
    return ret.flat();
  }
  function read_array_while(next) {
    let v3 = [];
    while (true) {
      let x2 = next(v3.length);
      if (!x2)
        break;
      v3.push(x2);
    }
    return v3;
  }
  function read_transposed(n4, w3, next) {
    let m4 = Array(n4).fill().map(() => []);
    for (let i4 = 0; i4 < w3; i4++) {
      read_deltas(n4, next).forEach((x2, j4) => m4[j4].push(x2));
    }
    return m4;
  }
  function read_linear_table(w3, next) {
    let dx = 1 + next();
    let dy = next();
    let vN = read_array_while(next);
    let m4 = read_transposed(vN.length, 1 + w3, next);
    return m4.flatMap((v3, i4) => {
      let [x2, ...ys] = v3;
      return Array(vN[i4]).fill().map((_, j4) => {
        let j_dy = j4 * dy;
        return [x2 + j4 * dx, ys.map((y2) => y2 + j_dy)];
      });
    });
  }
  function read_replacement_table(w3, next) {
    let n4 = 1 + next();
    let m4 = read_transposed(n4, 1 + w3, next);
    return m4.map((v3) => [v3[0], v3.slice(1)]);
  }
  var r$1 = read_compressed_payload("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");
  var FENCED = /* @__PURE__ */ new Map([[8217, "apostrophe"], [8260, "fraction slash"], [12539, "middle dot"]]);
  var NSM_MAX = 4;
  function hex_cp(cp) {
    return cp.toString(16).toUpperCase().padStart(2, "0");
  }
  function quote_cp(cp) {
    return `{${hex_cp(cp)}}`;
  }
  function explode_cp(s4) {
    let cps = [];
    for (let pos = 0, len = s4.length; pos < len; ) {
      let cp = s4.codePointAt(pos);
      pos += cp < 65536 ? 1 : 2;
      cps.push(cp);
    }
    return cps;
  }
  function str_from_cps(cps) {
    const chunk = 4096;
    let len = cps.length;
    if (len < chunk)
      return String.fromCodePoint(...cps);
    let buf = [];
    for (let i4 = 0; i4 < len; ) {
      buf.push(String.fromCodePoint(...cps.slice(i4, i4 += chunk)));
    }
    return buf.join("");
  }
  var r3 = read_compressed_payload("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");
  function unpack_cc(packed) {
    return packed >> 24 & 255;
  }
  function unpack_cp(packed) {
    return packed & 16777215;
  }
  var SHIFTED_RANK = new Map(read_sorted_arrays(r3).flatMap((v3, i4) => v3.map((x2) => [x2, i4 + 1 << 24])));
  var EXCLUSIONS = new Set(read_sorted(r3));
  var DECOMP = /* @__PURE__ */ new Map();
  var RECOMP = /* @__PURE__ */ new Map();
  for (let [cp, cps] of read_mapped(r3)) {
    if (!EXCLUSIONS.has(cp) && cps.length == 2) {
      let [a4, b5] = cps;
      let bucket = RECOMP.get(a4);
      if (!bucket) {
        bucket = /* @__PURE__ */ new Map();
        RECOMP.set(a4, bucket);
      }
      bucket.set(b5, cp);
    }
    DECOMP.set(cp, cps.reverse());
  }
  var S0 = 44032;
  var L0 = 4352;
  var V0 = 4449;
  var T0 = 4519;
  var L_COUNT = 19;
  var V_COUNT = 21;
  var T_COUNT = 28;
  var N_COUNT = V_COUNT * T_COUNT;
  var S_COUNT = L_COUNT * N_COUNT;
  var S1 = S0 + S_COUNT;
  var L1 = L0 + L_COUNT;
  var V1 = V0 + V_COUNT;
  var T1 = T0 + T_COUNT;
  function is_hangul(cp) {
    return cp >= S0 && cp < S1;
  }
  function compose_pair(a4, b5) {
    if (a4 >= L0 && a4 < L1 && b5 >= V0 && b5 < V1) {
      return S0 + (a4 - L0) * N_COUNT + (b5 - V0) * T_COUNT;
    } else if (is_hangul(a4) && b5 > T0 && b5 < T1 && (a4 - S0) % T_COUNT == 0) {
      return a4 + (b5 - T0);
    } else {
      let recomp = RECOMP.get(a4);
      if (recomp) {
        recomp = recomp.get(b5);
        if (recomp) {
          return recomp;
        }
      }
      return -1;
    }
  }
  function decomposed(cps) {
    let ret = [];
    let buf = [];
    let check_order = false;
    function add2(cp) {
      let cc = SHIFTED_RANK.get(cp);
      if (cc) {
        check_order = true;
        cp |= cc;
      }
      ret.push(cp);
    }
    for (let cp of cps) {
      while (true) {
        if (cp < 128) {
          ret.push(cp);
        } else if (is_hangul(cp)) {
          let s_index = cp - S0;
          let l_index = s_index / N_COUNT | 0;
          let v_index = s_index % N_COUNT / T_COUNT | 0;
          let t_index = s_index % T_COUNT;
          add2(L0 + l_index);
          add2(V0 + v_index);
          if (t_index > 0)
            add2(T0 + t_index);
        } else {
          let mapped = DECOMP.get(cp);
          if (mapped) {
            buf.push(...mapped);
          } else {
            add2(cp);
          }
        }
        if (!buf.length)
          break;
        cp = buf.pop();
      }
    }
    if (check_order && ret.length > 1) {
      let prev_cc = unpack_cc(ret[0]);
      for (let i4 = 1; i4 < ret.length; i4++) {
        let cc = unpack_cc(ret[i4]);
        if (cc == 0 || prev_cc <= cc) {
          prev_cc = cc;
          continue;
        }
        let j4 = i4 - 1;
        while (true) {
          let tmp = ret[j4 + 1];
          ret[j4 + 1] = ret[j4];
          ret[j4] = tmp;
          if (!j4)
            break;
          prev_cc = unpack_cc(ret[--j4]);
          if (prev_cc <= cc)
            break;
        }
        prev_cc = unpack_cc(ret[i4]);
      }
    }
    return ret;
  }
  function composed_from_decomposed(v3) {
    let ret = [];
    let stack = [];
    let prev_cp = -1;
    let prev_cc = 0;
    for (let packed of v3) {
      let cc = unpack_cc(packed);
      let cp = unpack_cp(packed);
      if (prev_cp == -1) {
        if (cc == 0) {
          prev_cp = cp;
        } else {
          ret.push(cp);
        }
      } else if (prev_cc > 0 && prev_cc >= cc) {
        if (cc == 0) {
          ret.push(prev_cp, ...stack);
          stack.length = 0;
          prev_cp = cp;
        } else {
          stack.push(cp);
        }
        prev_cc = cc;
      } else {
        let composed = compose_pair(prev_cp, cp);
        if (composed >= 0) {
          prev_cp = composed;
        } else if (prev_cc == 0 && cc == 0) {
          ret.push(prev_cp);
          prev_cp = cp;
        } else {
          stack.push(cp);
          prev_cc = cc;
        }
      }
    }
    if (prev_cp >= 0) {
      ret.push(prev_cp, ...stack);
    }
    return ret;
  }
  function nfd(cps) {
    return decomposed(cps).map(unpack_cp);
  }
  function nfc(cps) {
    return composed_from_decomposed(decomposed(cps));
  }
  var FE0F = 65039;
  var STOP_CH = ".";
  var UNIQUE_PH = 1;
  var HYPHEN = 45;
  function read_set() {
    return new Set(read_sorted(r$1));
  }
  var MAPPED = new Map(read_mapped(r$1));
  var IGNORED = read_set();
  var CM = read_set();
  var NSM = new Set(read_sorted(r$1).map(function(i4) {
    return this[i4];
  }, [...CM]));
  var ESCAPE = read_set();
  var NFC_CHECK = read_set();
  var CHUNKS = read_sorted_arrays(r$1);
  function read_chunked() {
    return new Set([read_sorted(r$1).map((i4) => CHUNKS[i4]), read_sorted(r$1)].flat(2));
  }
  var UNRESTRICTED = r$1();
  var GROUPS = read_array_while((i4) => {
    let N2 = read_array_while(r$1).map((x2) => x2 + 96);
    if (N2.length) {
      let R = i4 >= UNRESTRICTED;
      N2[0] -= 32;
      N2 = str_from_cps(N2);
      if (R)
        N2 = `Restricted[${N2}]`;
      let P2 = read_chunked();
      let Q = read_chunked();
      let V = [...P2, ...Q].sort((a4, b5) => a4 - b5);
      let M2 = !r$1();
      return { N: N2, P: P2, M: M2, R, V: new Set(V) };
    }
  });
  var WHOLE_VALID = read_set();
  var WHOLE_MAP = /* @__PURE__ */ new Map();
  [...WHOLE_VALID, ...read_set()].sort((a4, b5) => a4 - b5).map((cp, i4, v3) => {
    let d4 = r$1();
    let w3 = v3[i4] = d4 ? v3[i4 - d4] : { V: [], M: /* @__PURE__ */ new Map() };
    w3.V.push(cp);
    if (!WHOLE_VALID.has(cp)) {
      WHOLE_MAP.set(cp, w3);
    }
  });
  for (let { V, M: M2 } of new Set(WHOLE_MAP.values())) {
    let recs = [];
    for (let cp of V) {
      let gs = GROUPS.filter((g4) => g4.V.has(cp));
      let rec = recs.find(({ G }) => gs.some((g4) => G.has(g4)));
      if (!rec) {
        rec = { G: /* @__PURE__ */ new Set(), V: [] };
        recs.push(rec);
      }
      rec.V.push(cp);
      gs.forEach((g4) => rec.G.add(g4));
    }
    let union2 = recs.flatMap(({ G }) => [...G]);
    for (let { G, V: V2 } of recs) {
      let complement = new Set(union2.filter((g4) => !G.has(g4)));
      for (let cp of V2) {
        M2.set(cp, complement);
      }
    }
  }
  var union = /* @__PURE__ */ new Set();
  var multi = /* @__PURE__ */ new Set();
  for (let g4 of GROUPS) {
    for (let cp of g4.V) {
      (union.has(cp) ? multi : union).add(cp);
    }
  }
  for (let cp of union) {
    if (!WHOLE_MAP.has(cp) && !multi.has(cp)) {
      WHOLE_MAP.set(cp, UNIQUE_PH);
    }
  }
  var VALID = /* @__PURE__ */ new Set([...union, ...nfd(union)]);
  var EMOJI_SORTED = read_sorted(r$1);
  var EMOJI_ROOT = read_emoji_trie([]);
  function read_emoji_trie(cps) {
    let B3 = read_array_while(() => {
      let keys = read_sorted(r$1).map((i4) => EMOJI_SORTED[i4]);
      if (keys.length)
        return read_emoji_trie(keys);
    }).sort((a4, b5) => b5.Q.size - a4.Q.size);
    let temp = r$1();
    let V = temp % 3;
    temp = temp / 3 | 0;
    let F2 = temp & 1;
    temp >>= 1;
    let S2 = temp & 1;
    let C = temp & 2;
    return { B: B3, V, F: F2, S: S2, C, Q: new Set(cps) };
  }
  var Emoji = class extends Array {
    get is_emoji() {
      return true;
    }
  };
  function safe_str_from_cps(cps, quoter = quote_cp) {
    let buf = [];
    if (is_combining_mark(cps[0]))
      buf.push("\u25CC");
    let prev = 0;
    let n4 = cps.length;
    for (let i4 = 0; i4 < n4; i4++) {
      let cp = cps[i4];
      if (should_escape(cp)) {
        buf.push(str_from_cps(cps.slice(prev, i4)));
        buf.push(quoter(cp));
        prev = i4 + 1;
      }
    }
    buf.push(str_from_cps(cps.slice(prev, n4)));
    return buf.join("");
  }
  function quoted_cp(cp) {
    return (should_escape(cp) ? "" : `${bidi_qq(safe_str_from_cps([cp]))} `) + quote_cp(cp);
  }
  function bidi_qq(s4) {
    return `"${s4}"\u200E`;
  }
  function check_label_extension(cps) {
    if (cps.length >= 4 && cps[2] == HYPHEN && cps[3] == HYPHEN) {
      throw new Error("invalid label extension");
    }
  }
  function check_leading_underscore(cps) {
    const UNDERSCORE = 95;
    for (let i4 = cps.lastIndexOf(UNDERSCORE); i4 > 0; ) {
      if (cps[--i4] !== UNDERSCORE) {
        throw new Error("underscore allowed only at start");
      }
    }
  }
  function check_fenced(cps) {
    let cp = cps[0];
    let prev = FENCED.get(cp);
    if (prev)
      throw error_placement(`leading ${prev}`);
    let n4 = cps.length;
    let last = -1;
    for (let i4 = 1; i4 < n4; i4++) {
      cp = cps[i4];
      let match = FENCED.get(cp);
      if (match) {
        if (last == i4)
          throw error_placement(`${prev} + ${match}`);
        last = i4 + 1;
        prev = match;
      }
    }
    if (last == n4)
      throw error_placement(`trailing ${prev}`);
  }
  function is_combining_mark(cp) {
    return CM.has(cp);
  }
  function should_escape(cp) {
    return ESCAPE.has(cp);
  }
  function ens_normalize(name) {
    return flatten(ens_split(name));
  }
  function ens_split(name, preserve_emoji) {
    let offset = 0;
    return name.split(STOP_CH).map((label) => {
      let input = explode_cp(label);
      let info = {
        input,
        offset
        // codepoint, not substring!
      };
      offset += input.length + 1;
      let norm;
      try {
        let tokens = info.tokens = process(input, nfc);
        let token_count = tokens.length;
        let type;
        if (!token_count) {
          throw new Error(`empty label`);
        } else {
          let chars = tokens[0];
          let emoji = token_count > 1 || chars.is_emoji;
          if (!emoji && chars.every((cp) => cp < 128)) {
            norm = chars;
            check_leading_underscore(norm);
            check_label_extension(norm);
            type = "ASCII";
          } else {
            if (emoji) {
              info.emoji = true;
              chars = tokens.flatMap((x2) => x2.is_emoji ? [] : x2);
            }
            norm = tokens.flatMap((x2) => !preserve_emoji && x2.is_emoji ? filter_fe0f(x2) : x2);
            check_leading_underscore(norm);
            if (!chars.length) {
              type = "Emoji";
            } else {
              if (CM.has(norm[0]))
                throw error_placement("leading combining mark");
              for (let i4 = 1; i4 < token_count; i4++) {
                let cps = tokens[i4];
                if (!cps.is_emoji && CM.has(cps[0])) {
                  throw error_placement(`emoji + combining mark: "${str_from_cps(tokens[i4 - 1])} + ${safe_str_from_cps([cps[0]])}"`);
                }
              }
              check_fenced(norm);
              let unique = [...new Set(chars)];
              let [g4] = determine_group(unique);
              check_group(g4, chars);
              check_whole(g4, unique);
              type = g4.N;
            }
          }
        }
        info.type = type;
      } catch (err) {
        info.error = err;
      }
      info.output = norm;
      return info;
    });
  }
  function check_whole(group, unique) {
    let maker;
    let shared = [];
    for (let cp of unique) {
      let whole = WHOLE_MAP.get(cp);
      if (whole === UNIQUE_PH)
        return;
      if (whole) {
        let set = whole.M.get(cp);
        maker = maker ? maker.filter((g4) => set.has(g4)) : [...set];
        if (!maker.length)
          return;
      } else {
        shared.push(cp);
      }
    }
    if (maker) {
      for (let g4 of maker) {
        if (shared.every((cp) => g4.V.has(cp))) {
          throw new Error(`whole-script confusable: ${group.N}/${g4.N}`);
        }
      }
    }
  }
  function determine_group(unique) {
    let groups = GROUPS;
    for (let cp of unique) {
      let gs = groups.filter((g4) => g4.V.has(cp));
      if (!gs.length) {
        if (groups === GROUPS) {
          throw error_disallowed(cp);
        } else {
          throw error_group_member(groups[0], cp);
        }
      }
      groups = gs;
      if (gs.length == 1)
        break;
    }
    return groups;
  }
  function flatten(split2) {
    return split2.map(({ input, error, output: output2 }) => {
      if (error) {
        let msg = error.message;
        throw new Error(split2.length == 1 ? msg : `Invalid label ${bidi_qq(safe_str_from_cps(input))}: ${msg}`);
      }
      return str_from_cps(output2);
    }).join(STOP_CH);
  }
  function error_disallowed(cp) {
    return new Error(`disallowed character: ${quoted_cp(cp)}`);
  }
  function error_group_member(g4, cp) {
    let quoted = quoted_cp(cp);
    let gg = GROUPS.find((g5) => g5.P.has(cp));
    if (gg) {
      quoted = `${gg.N} ${quoted}`;
    }
    return new Error(`illegal mixture: ${g4.N} + ${quoted}`);
  }
  function error_placement(where) {
    return new Error(`illegal placement: ${where}`);
  }
  function check_group(g4, cps) {
    let { V, M: M2 } = g4;
    for (let cp of cps) {
      if (!V.has(cp)) {
        throw error_group_member(g4, cp);
      }
    }
    if (M2) {
      let decomposed2 = nfd(cps);
      for (let i4 = 1, e4 = decomposed2.length; i4 < e4; i4++) {
        if (NSM.has(decomposed2[i4])) {
          let j4 = i4 + 1;
          for (let cp; j4 < e4 && NSM.has(cp = decomposed2[j4]); j4++) {
            for (let k3 = i4; k3 < j4; k3++) {
              if (decomposed2[k3] == cp) {
                throw new Error(`non-spacing marks: repeated ${quoted_cp(cp)}`);
              }
            }
          }
          if (j4 - i4 > NSM_MAX) {
            throw new Error(`non-spacing marks: too many ${bidi_qq(safe_str_from_cps(decomposed2.slice(i4 - 1, j4)))} (${j4 - i4}/${NSM_MAX})`);
          }
          i4 = j4;
        }
      }
    }
  }
  function process(input, nf) {
    let ret = [];
    let chars = [];
    input = input.slice().reverse();
    while (input.length) {
      let emoji = consume_emoji_reversed(input);
      if (emoji) {
        if (chars.length) {
          ret.push(nf(chars));
          chars = [];
        }
        ret.push(emoji);
      } else {
        let cp = input.pop();
        if (VALID.has(cp)) {
          chars.push(cp);
        } else {
          let cps = MAPPED.get(cp);
          if (cps) {
            chars.push(...cps);
          } else if (!IGNORED.has(cp)) {
            throw error_disallowed(cp);
          }
        }
      }
    }
    if (chars.length) {
      ret.push(nf(chars));
    }
    return ret;
  }
  function filter_fe0f(cps) {
    return cps.filter((cp) => cp != FE0F);
  }
  function consume_emoji_reversed(cps, eaten) {
    let node = EMOJI_ROOT;
    let emoji;
    let saved;
    let stack = [];
    let pos = cps.length;
    if (eaten)
      eaten.length = 0;
    while (pos) {
      let cp = cps[--pos];
      node = node.B.find((x2) => x2.Q.has(cp));
      if (!node)
        break;
      if (node.S) {
        saved = cp;
      } else if (node.C) {
        if (cp === saved)
          break;
      }
      stack.push(cp);
      if (node.F) {
        stack.push(FE0F);
        if (pos > 0 && cps[pos - 1] == FE0F)
          pos--;
      }
      if (node.V) {
        emoji = conform_emoji_copy(stack, node);
        if (eaten)
          eaten.push(...cps.slice(pos).reverse());
        cps.length = pos;
      }
    }
    return emoji;
  }
  function conform_emoji_copy(cps, node) {
    let copy4 = Emoji.from(cps);
    if (node.V == 2)
      copy4.splice(1, 1);
    return copy4;
  }

  // node_modules/ethers/lib.esm/hash/namehash.js
  var Zeros2 = new Uint8Array(32);
  Zeros2.fill(0);
  function checkComponent(comp) {
    assertArgument(comp.length !== 0, "invalid ENS name; empty component", "comp", comp);
    return comp;
  }
  function ensNameSplit(name) {
    const bytes2 = toUtf8Bytes(ensNormalize(name));
    const comps = [];
    if (name.length === 0) {
      return comps;
    }
    let last = 0;
    for (let i4 = 0; i4 < bytes2.length; i4++) {
      const d4 = bytes2[i4];
      if (d4 === 46) {
        comps.push(checkComponent(bytes2.slice(last, i4)));
        last = i4 + 1;
      }
    }
    assertArgument(last < bytes2.length, "invalid ENS name; empty component", "name", name);
    comps.push(checkComponent(bytes2.slice(last)));
    return comps;
  }
  function ensNormalize(name) {
    try {
      return ens_normalize(name);
    } catch (error) {
      assertArgument(false, `invalid ENS name (${error.message})`, "name", name);
    }
  }
  function namehash(name) {
    assertArgument(typeof name === "string", "invalid ENS name; not a string", "name", name);
    let result = Zeros2;
    const comps = ensNameSplit(name);
    while (comps.length) {
      result = keccak256(concat([result, keccak256(comps.pop())]));
    }
    return hexlify(result);
  }
  function dnsEncode(name) {
    return hexlify(concat(ensNameSplit(name).map((comp) => {
      if (comp.length > 63) {
        throw new Error("invalid DNS encoded entry; length exceeds 63 bytes");
      }
      const bytes2 = new Uint8Array(comp.length + 1);
      bytes2.set(comp, 1);
      bytes2[0] = bytes2.length - 1;
      return bytes2;
    }))) + "00";
  }

  // node_modules/ethers/lib.esm/transaction/accesslist.js
  function accessSetify(addr, storageKeys) {
    return {
      address: getAddress(addr),
      storageKeys: storageKeys.map((storageKey, index) => {
        assertArgument(isHexString(storageKey, 32), "invalid slot", `storageKeys[${index}]`, storageKey);
        return storageKey.toLowerCase();
      })
    };
  }
  function accessListify(value) {
    if (Array.isArray(value)) {
      return value.map((set, index) => {
        if (Array.isArray(set)) {
          assertArgument(set.length === 2, "invalid slot set", `value[${index}]`, set);
          return accessSetify(set[0], set[1]);
        }
        assertArgument(set != null && typeof set === "object", "invalid address-slot set", "value", value);
        return accessSetify(set.address, set.storageKeys);
      });
    }
    assertArgument(value != null && typeof value === "object", "invalid access list", "value", value);
    const result = Object.keys(value).map((addr) => {
      const storageKeys = value[addr].reduce((accum, storageKey) => {
        accum[storageKey] = true;
        return accum;
      }, {});
      return accessSetify(addr, Object.keys(storageKeys).sort());
    });
    result.sort((a4, b5) => a4.address.localeCompare(b5.address));
    return result;
  }

  // node_modules/ethers/lib.esm/transaction/address.js
  function computeAddress(key) {
    let pubkey;
    if (typeof key === "string") {
      pubkey = SigningKey.computePublicKey(key, false);
    } else {
      pubkey = key.publicKey;
    }
    return getAddress(keccak256("0x" + pubkey.substring(4)).substring(26));
  }
  function recoverAddress(digest, signature) {
    return computeAddress(SigningKey.recoverPublicKey(digest, signature));
  }

  // node_modules/ethers/lib.esm/transaction/transaction.js
  var BN_07 = BigInt(0);
  var BN_22 = BigInt(2);
  var BN_272 = BigInt(27);
  var BN_282 = BigInt(28);
  var BN_352 = BigInt(35);
  var BN_MAX_UINT = BigInt("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
  function handleAddress(value) {
    if (value === "0x") {
      return null;
    }
    return getAddress(value);
  }
  function handleAccessList(value, param) {
    try {
      return accessListify(value);
    } catch (error) {
      assertArgument(false, error.message, param, value);
    }
  }
  function handleNumber(_value, param) {
    if (_value === "0x") {
      return 0;
    }
    return getNumber(_value, param);
  }
  function handleUint(_value, param) {
    if (_value === "0x") {
      return BN_07;
    }
    const value = getBigInt(_value, param);
    assertArgument(value <= BN_MAX_UINT, "value exceeds uint size", param, value);
    return value;
  }
  function formatNumber(_value, name) {
    const value = getBigInt(_value, "value");
    const result = toBeArray(value);
    assertArgument(result.length <= 32, `value too large`, `tx.${name}`, value);
    return result;
  }
  function formatAccessList(value) {
    return accessListify(value).map((set) => [set.address, set.storageKeys]);
  }
  function _parseLegacy(data) {
    const fields = decodeRlp(data);
    assertArgument(Array.isArray(fields) && (fields.length === 9 || fields.length === 6), "invalid field count for legacy transaction", "data", data);
    const tx = {
      type: 0,
      nonce: handleNumber(fields[0], "nonce"),
      gasPrice: handleUint(fields[1], "gasPrice"),
      gasLimit: handleUint(fields[2], "gasLimit"),
      to: handleAddress(fields[3]),
      value: handleUint(fields[4], "value"),
      data: hexlify(fields[5]),
      chainId: BN_07
    };
    if (fields.length === 6) {
      return tx;
    }
    const v3 = handleUint(fields[6], "v");
    const r4 = handleUint(fields[7], "r");
    const s4 = handleUint(fields[8], "s");
    if (r4 === BN_07 && s4 === BN_07) {
      tx.chainId = v3;
    } else {
      let chainId = (v3 - BN_352) / BN_22;
      if (chainId < BN_07) {
        chainId = BN_07;
      }
      tx.chainId = chainId;
      assertArgument(chainId !== BN_07 || (v3 === BN_272 || v3 === BN_282), "non-canonical legacy v", "v", fields[6]);
      tx.signature = Signature2.from({
        r: zeroPadValue(fields[7], 32),
        s: zeroPadValue(fields[8], 32),
        v: v3
      });
      tx.hash = keccak256(data);
    }
    return tx;
  }
  function _serializeLegacy(tx, sig) {
    const fields = [
      formatNumber(tx.nonce || 0, "nonce"),
      formatNumber(tx.gasPrice || 0, "gasPrice"),
      formatNumber(tx.gasLimit || 0, "gasLimit"),
      tx.to != null ? getAddress(tx.to) : "0x",
      formatNumber(tx.value || 0, "value"),
      tx.data || "0x"
    ];
    let chainId = BN_07;
    if (tx.chainId != BN_07) {
      chainId = getBigInt(tx.chainId, "tx.chainId");
      assertArgument(!sig || sig.networkV == null || sig.legacyChainId === chainId, "tx.chainId/sig.v mismatch", "sig", sig);
    } else if (tx.signature) {
      const legacy = tx.signature.legacyChainId;
      if (legacy != null) {
        chainId = legacy;
      }
    }
    if (!sig) {
      if (chainId !== BN_07) {
        fields.push(toBeArray(chainId));
        fields.push("0x");
        fields.push("0x");
      }
      return encodeRlp(fields);
    }
    let v3 = BigInt(27 + sig.yParity);
    if (chainId !== BN_07) {
      v3 = Signature2.getChainIdV(chainId, sig.v);
    } else if (BigInt(sig.v) !== v3) {
      assertArgument(false, "tx.chainId/sig.v mismatch", "sig", sig);
    }
    fields.push(toBeArray(v3));
    fields.push(toBeArray(sig.r));
    fields.push(toBeArray(sig.s));
    return encodeRlp(fields);
  }
  function _parseEipSignature(tx, fields) {
    let yParity;
    try {
      yParity = handleNumber(fields[0], "yParity");
      if (yParity !== 0 && yParity !== 1) {
        throw new Error("bad yParity");
      }
    } catch (error) {
      assertArgument(false, "invalid yParity", "yParity", fields[0]);
    }
    const r4 = zeroPadValue(fields[1], 32);
    const s4 = zeroPadValue(fields[2], 32);
    const signature = Signature2.from({ r: r4, s: s4, yParity });
    tx.signature = signature;
  }
  function _parseEip1559(data) {
    const fields = decodeRlp(getBytes(data).slice(1));
    assertArgument(Array.isArray(fields) && (fields.length === 9 || fields.length === 12), "invalid field count for transaction type: 2", "data", hexlify(data));
    const maxPriorityFeePerGas = handleUint(fields[2], "maxPriorityFeePerGas");
    const maxFeePerGas = handleUint(fields[3], "maxFeePerGas");
    const tx = {
      type: 2,
      chainId: handleUint(fields[0], "chainId"),
      nonce: handleNumber(fields[1], "nonce"),
      maxPriorityFeePerGas,
      maxFeePerGas,
      gasPrice: null,
      gasLimit: handleUint(fields[4], "gasLimit"),
      to: handleAddress(fields[5]),
      value: handleUint(fields[6], "value"),
      data: hexlify(fields[7]),
      accessList: handleAccessList(fields[8], "accessList")
    };
    if (fields.length === 9) {
      return tx;
    }
    tx.hash = keccak256(data);
    _parseEipSignature(tx, fields.slice(9));
    return tx;
  }
  function _serializeEip1559(tx, sig) {
    const fields = [
      formatNumber(tx.chainId || 0, "chainId"),
      formatNumber(tx.nonce || 0, "nonce"),
      formatNumber(tx.maxPriorityFeePerGas || 0, "maxPriorityFeePerGas"),
      formatNumber(tx.maxFeePerGas || 0, "maxFeePerGas"),
      formatNumber(tx.gasLimit || 0, "gasLimit"),
      tx.to != null ? getAddress(tx.to) : "0x",
      formatNumber(tx.value || 0, "value"),
      tx.data || "0x",
      formatAccessList(tx.accessList || [])
    ];
    if (sig) {
      fields.push(formatNumber(sig.yParity, "yParity"));
      fields.push(toBeArray(sig.r));
      fields.push(toBeArray(sig.s));
    }
    return concat(["0x02", encodeRlp(fields)]);
  }
  function _parseEip2930(data) {
    const fields = decodeRlp(getBytes(data).slice(1));
    assertArgument(Array.isArray(fields) && (fields.length === 8 || fields.length === 11), "invalid field count for transaction type: 1", "data", hexlify(data));
    const tx = {
      type: 1,
      chainId: handleUint(fields[0], "chainId"),
      nonce: handleNumber(fields[1], "nonce"),
      gasPrice: handleUint(fields[2], "gasPrice"),
      gasLimit: handleUint(fields[3], "gasLimit"),
      to: handleAddress(fields[4]),
      value: handleUint(fields[5], "value"),
      data: hexlify(fields[6]),
      accessList: handleAccessList(fields[7], "accessList")
    };
    if (fields.length === 8) {
      return tx;
    }
    tx.hash = keccak256(data);
    _parseEipSignature(tx, fields.slice(8));
    return tx;
  }
  function _serializeEip2930(tx, sig) {
    const fields = [
      formatNumber(tx.chainId || 0, "chainId"),
      formatNumber(tx.nonce || 0, "nonce"),
      formatNumber(tx.gasPrice || 0, "gasPrice"),
      formatNumber(tx.gasLimit || 0, "gasLimit"),
      tx.to != null ? getAddress(tx.to) : "0x",
      formatNumber(tx.value || 0, "value"),
      tx.data || "0x",
      formatAccessList(tx.accessList || [])
    ];
    if (sig) {
      fields.push(formatNumber(sig.yParity, "recoveryParam"));
      fields.push(toBeArray(sig.r));
      fields.push(toBeArray(sig.s));
    }
    return concat(["0x01", encodeRlp(fields)]);
  }
  var Transaction = class _Transaction {
    #type;
    #to;
    #data;
    #nonce;
    #gasLimit;
    #gasPrice;
    #maxPriorityFeePerGas;
    #maxFeePerGas;
    #value;
    #chainId;
    #sig;
    #accessList;
    /**
     *  The transaction type.
     *
     *  If null, the type will be automatically inferred based on
     *  explicit properties.
     */
    get type() {
      return this.#type;
    }
    set type(value) {
      switch (value) {
        case null:
          this.#type = null;
          break;
        case 0:
        case "legacy":
          this.#type = 0;
          break;
        case 1:
        case "berlin":
        case "eip-2930":
          this.#type = 1;
          break;
        case 2:
        case "london":
        case "eip-1559":
          this.#type = 2;
          break;
        default:
          assertArgument(false, "unsupported transaction type", "type", value);
      }
    }
    /**
     *  The name of the transaction type.
     */
    get typeName() {
      switch (this.type) {
        case 0:
          return "legacy";
        case 1:
          return "eip-2930";
        case 2:
          return "eip-1559";
      }
      return null;
    }
    /**
     *  The ``to`` address for the transaction or ``null`` if the
     *  transaction is an ``init`` transaction.
     */
    get to() {
      return this.#to;
    }
    set to(value) {
      this.#to = value == null ? null : getAddress(value);
    }
    /**
     *  The transaction nonce.
     */
    get nonce() {
      return this.#nonce;
    }
    set nonce(value) {
      this.#nonce = getNumber(value, "value");
    }
    /**
     *  The gas limit.
     */
    get gasLimit() {
      return this.#gasLimit;
    }
    set gasLimit(value) {
      this.#gasLimit = getBigInt(value);
    }
    /**
     *  The gas price.
     *
     *  On legacy networks this defines the fee that will be paid. On
     *  EIP-1559 networks, this should be ``null``.
     */
    get gasPrice() {
      const value = this.#gasPrice;
      if (value == null && (this.type === 0 || this.type === 1)) {
        return BN_07;
      }
      return value;
    }
    set gasPrice(value) {
      this.#gasPrice = value == null ? null : getBigInt(value, "gasPrice");
    }
    /**
     *  The maximum priority fee per unit of gas to pay. On legacy
     *  networks this should be ``null``.
     */
    get maxPriorityFeePerGas() {
      const value = this.#maxPriorityFeePerGas;
      if (value == null) {
        if (this.type === 2) {
          return BN_07;
        }
        return null;
      }
      return value;
    }
    set maxPriorityFeePerGas(value) {
      this.#maxPriorityFeePerGas = value == null ? null : getBigInt(value, "maxPriorityFeePerGas");
    }
    /**
     *  The maximum total fee per unit of gas to pay. On legacy
     *  networks this should be ``null``.
     */
    get maxFeePerGas() {
      const value = this.#maxFeePerGas;
      if (value == null) {
        if (this.type === 2) {
          return BN_07;
        }
        return null;
      }
      return value;
    }
    set maxFeePerGas(value) {
      this.#maxFeePerGas = value == null ? null : getBigInt(value, "maxFeePerGas");
    }
    /**
     *  The transaction data. For ``init`` transactions this is the
     *  deployment code.
     */
    get data() {
      return this.#data;
    }
    set data(value) {
      this.#data = hexlify(value);
    }
    /**
     *  The amount of ether (in wei) to send in this transactions.
     */
    get value() {
      return this.#value;
    }
    set value(value) {
      this.#value = getBigInt(value, "value");
    }
    /**
     *  The chain ID this transaction is valid on.
     */
    get chainId() {
      return this.#chainId;
    }
    set chainId(value) {
      this.#chainId = getBigInt(value);
    }
    /**
     *  If signed, the signature for this transaction.
     */
    get signature() {
      return this.#sig || null;
    }
    set signature(value) {
      this.#sig = value == null ? null : Signature2.from(value);
    }
    /**
     *  The access list.
     *
     *  An access list permits discounted (but pre-paid) access to
     *  bytecode and state variable access within contract execution.
     */
    get accessList() {
      const value = this.#accessList || null;
      if (value == null) {
        if (this.type === 1 || this.type === 2) {
          return [];
        }
        return null;
      }
      return value;
    }
    set accessList(value) {
      this.#accessList = value == null ? null : accessListify(value);
    }
    /**
     *  Creates a new Transaction with default values.
     */
    constructor() {
      this.#type = null;
      this.#to = null;
      this.#nonce = 0;
      this.#gasLimit = BigInt(0);
      this.#gasPrice = null;
      this.#maxPriorityFeePerGas = null;
      this.#maxFeePerGas = null;
      this.#data = "0x";
      this.#value = BigInt(0);
      this.#chainId = BigInt(0);
      this.#sig = null;
      this.#accessList = null;
    }
    /**
     *  The transaction hash, if signed. Otherwise, ``null``.
     */
    get hash() {
      if (this.signature == null) {
        return null;
      }
      return keccak256(this.serialized);
    }
    /**
     *  The pre-image hash of this transaction.
     *
     *  This is the digest that a [[Signer]] must sign to authorize
     *  this transaction.
     */
    get unsignedHash() {
      return keccak256(this.unsignedSerialized);
    }
    /**
     *  The sending address, if signed. Otherwise, ``null``.
     */
    get from() {
      if (this.signature == null) {
        return null;
      }
      return recoverAddress(this.unsignedHash, this.signature);
    }
    /**
     *  The public key of the sender, if signed. Otherwise, ``null``.
     */
    get fromPublicKey() {
      if (this.signature == null) {
        return null;
      }
      return SigningKey.recoverPublicKey(this.unsignedHash, this.signature);
    }
    /**
     *  Returns true if signed.
     *
     *  This provides a Type Guard that properties requiring a signed
     *  transaction are non-null.
     */
    isSigned() {
      return this.signature != null;
    }
    /**
     *  The serialized transaction.
     *
     *  This throws if the transaction is unsigned. For the pre-image,
     *  use [[unsignedSerialized]].
     */
    get serialized() {
      assert(this.signature != null, "cannot serialize unsigned transaction; maybe you meant .unsignedSerialized", "UNSUPPORTED_OPERATION", { operation: ".serialized" });
      switch (this.inferType()) {
        case 0:
          return _serializeLegacy(this, this.signature);
        case 1:
          return _serializeEip2930(this, this.signature);
        case 2:
          return _serializeEip1559(this, this.signature);
      }
      assert(false, "unsupported transaction type", "UNSUPPORTED_OPERATION", { operation: ".serialized" });
    }
    /**
     *  The transaction pre-image.
     *
     *  The hash of this is the digest which needs to be signed to
     *  authorize this transaction.
     */
    get unsignedSerialized() {
      switch (this.inferType()) {
        case 0:
          return _serializeLegacy(this);
        case 1:
          return _serializeEip2930(this);
        case 2:
          return _serializeEip1559(this);
      }
      assert(false, "unsupported transaction type", "UNSUPPORTED_OPERATION", { operation: ".unsignedSerialized" });
    }
    /**
     *  Return the most "likely" type; currently the highest
     *  supported transaction type.
     */
    inferType() {
      return this.inferTypes().pop();
    }
    /**
     *  Validates the explicit properties and returns a list of compatible
     *  transaction types.
     */
    inferTypes() {
      const hasGasPrice = this.gasPrice != null;
      const hasFee = this.maxFeePerGas != null || this.maxPriorityFeePerGas != null;
      const hasAccessList = this.accessList != null;
      if (this.maxFeePerGas != null && this.maxPriorityFeePerGas != null) {
        assert(this.maxFeePerGas >= this.maxPriorityFeePerGas, "priorityFee cannot be more than maxFee", "BAD_DATA", { value: this });
      }
      assert(!hasFee || this.type !== 0 && this.type !== 1, "transaction type cannot have maxFeePerGas or maxPriorityFeePerGas", "BAD_DATA", { value: this });
      assert(this.type !== 0 || !hasAccessList, "legacy transaction cannot have accessList", "BAD_DATA", { value: this });
      const types = [];
      if (this.type != null) {
        types.push(this.type);
      } else {
        if (hasFee) {
          types.push(2);
        } else if (hasGasPrice) {
          types.push(1);
          if (!hasAccessList) {
            types.push(0);
          }
        } else if (hasAccessList) {
          types.push(1);
          types.push(2);
        } else {
          types.push(0);
          types.push(1);
          types.push(2);
        }
      }
      types.sort();
      return types;
    }
    /**
     *  Returns true if this transaction is a legacy transaction (i.e.
     *  ``type === 0``).
     *
     *  This provides a Type Guard that the related properties are
     *  non-null.
     */
    isLegacy() {
      return this.type === 0;
    }
    /**
     *  Returns true if this transaction is berlin hardform transaction (i.e.
     *  ``type === 1``).
     *
     *  This provides a Type Guard that the related properties are
     *  non-null.
     */
    isBerlin() {
      return this.type === 1;
    }
    /**
     *  Returns true if this transaction is london hardform transaction (i.e.
     *  ``type === 2``).
     *
     *  This provides a Type Guard that the related properties are
     *  non-null.
     */
    isLondon() {
      return this.type === 2;
    }
    /**
     *  Create a copy of this transaciton.
     */
    clone() {
      return _Transaction.from(this);
    }
    /**
     *  Return a JSON-friendly object.
     */
    toJSON() {
      const s4 = (v3) => {
        if (v3 == null) {
          return null;
        }
        return v3.toString();
      };
      return {
        type: this.type,
        to: this.to,
        //            from: this.from,
        data: this.data,
        nonce: this.nonce,
        gasLimit: s4(this.gasLimit),
        gasPrice: s4(this.gasPrice),
        maxPriorityFeePerGas: s4(this.maxPriorityFeePerGas),
        maxFeePerGas: s4(this.maxFeePerGas),
        value: s4(this.value),
        chainId: s4(this.chainId),
        sig: this.signature ? this.signature.toJSON() : null,
        accessList: this.accessList
      };
    }
    /**
     *  Create a **Transaction** from a serialized transaction or a
     *  Transaction-like object.
     */
    static from(tx) {
      if (tx == null) {
        return new _Transaction();
      }
      if (typeof tx === "string") {
        const payload = getBytes(tx);
        if (payload[0] >= 127) {
          return _Transaction.from(_parseLegacy(payload));
        }
        switch (payload[0]) {
          case 1:
            return _Transaction.from(_parseEip2930(payload));
          case 2:
            return _Transaction.from(_parseEip1559(payload));
        }
        assert(false, "unsupported transaction type", "UNSUPPORTED_OPERATION", { operation: "from" });
      }
      const result = new _Transaction();
      if (tx.type != null) {
        result.type = tx.type;
      }
      if (tx.to != null) {
        result.to = tx.to;
      }
      if (tx.nonce != null) {
        result.nonce = tx.nonce;
      }
      if (tx.gasLimit != null) {
        result.gasLimit = tx.gasLimit;
      }
      if (tx.gasPrice != null) {
        result.gasPrice = tx.gasPrice;
      }
      if (tx.maxPriorityFeePerGas != null) {
        result.maxPriorityFeePerGas = tx.maxPriorityFeePerGas;
      }
      if (tx.maxFeePerGas != null) {
        result.maxFeePerGas = tx.maxFeePerGas;
      }
      if (tx.data != null) {
        result.data = tx.data;
      }
      if (tx.value != null) {
        result.value = tx.value;
      }
      if (tx.chainId != null) {
        result.chainId = tx.chainId;
      }
      if (tx.signature != null) {
        result.signature = Signature2.from(tx.signature);
      }
      if (tx.accessList != null) {
        result.accessList = tx.accessList;
      }
      if (tx.hash != null) {
        assertArgument(result.isSigned(), "unsigned transaction cannot define hash", "tx", tx);
        assertArgument(result.hash === tx.hash, "hash mismatch", "tx", tx);
      }
      if (tx.from != null) {
        assertArgument(result.isSigned(), "unsigned transaction cannot define from", "tx", tx);
        assertArgument(result.from.toLowerCase() === (tx.from || "").toLowerCase(), "from mismatch", "tx", tx);
      }
      return result;
    }
  };

  // node_modules/ethers/lib.esm/hash/typed-data.js
  var padding = new Uint8Array(32);
  padding.fill(0);
  var BN__1 = BigInt(-1);
  var BN_08 = BigInt(0);
  var BN_15 = BigInt(1);
  var BN_MAX_UINT2562 = BigInt("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff");
  function hexPadRight(value) {
    const bytes2 = getBytes(value);
    const padOffset = bytes2.length % 32;
    if (padOffset) {
      return concat([bytes2, padding.slice(padOffset)]);
    }
    return hexlify(bytes2);
  }
  var hexTrue = toBeHex(BN_15, 32);
  var hexFalse = toBeHex(BN_08, 32);
  var domainFieldTypes = {
    name: "string",
    version: "string",
    chainId: "uint256",
    verifyingContract: "address",
    salt: "bytes32"
  };
  var domainFieldNames = [
    "name",
    "version",
    "chainId",
    "verifyingContract",
    "salt"
  ];
  function checkString(key) {
    return function(value) {
      assertArgument(typeof value === "string", `invalid domain value for ${JSON.stringify(key)}`, `domain.${key}`, value);
      return value;
    };
  }
  var domainChecks = {
    name: checkString("name"),
    version: checkString("version"),
    chainId: function(_value) {
      const value = getBigInt(_value, "domain.chainId");
      assertArgument(value >= 0, "invalid chain ID", "domain.chainId", _value);
      if (Number.isSafeInteger(value)) {
        return Number(value);
      }
      return toQuantity(value);
    },
    verifyingContract: function(value) {
      try {
        return getAddress(value).toLowerCase();
      } catch (error) {
      }
      assertArgument(false, `invalid domain value "verifyingContract"`, "domain.verifyingContract", value);
    },
    salt: function(value) {
      const bytes2 = getBytes(value, "domain.salt");
      assertArgument(bytes2.length === 32, `invalid domain value "salt"`, "domain.salt", value);
      return hexlify(bytes2);
    }
  };
  function getBaseEncoder(type) {
    {
      const match = type.match(/^(u?)int(\d*)$/);
      if (match) {
        const signed2 = match[1] === "";
        const width = parseInt(match[2] || "256");
        assertArgument(width % 8 === 0 && width !== 0 && width <= 256 && (match[2] == null || match[2] === String(width)), "invalid numeric width", "type", type);
        const boundsUpper = mask(BN_MAX_UINT2562, signed2 ? width - 1 : width);
        const boundsLower = signed2 ? (boundsUpper + BN_15) * BN__1 : BN_08;
        return function(_value) {
          const value = getBigInt(_value, "value");
          assertArgument(value >= boundsLower && value <= boundsUpper, `value out-of-bounds for ${type}`, "value", value);
          return toBeHex(signed2 ? toTwos(value, 256) : value, 32);
        };
      }
    }
    {
      const match = type.match(/^bytes(\d+)$/);
      if (match) {
        const width = parseInt(match[1]);
        assertArgument(width !== 0 && width <= 32 && match[1] === String(width), "invalid bytes width", "type", type);
        return function(value) {
          const bytes2 = getBytes(value);
          assertArgument(bytes2.length === width, `invalid length for ${type}`, "value", value);
          return hexPadRight(value);
        };
      }
    }
    switch (type) {
      case "address":
        return function(value) {
          return zeroPadValue(getAddress(value), 32);
        };
      case "bool":
        return function(value) {
          return !value ? hexFalse : hexTrue;
        };
      case "bytes":
        return function(value) {
          return keccak256(value);
        };
      case "string":
        return function(value) {
          return id(value);
        };
    }
    return null;
  }
  function encodeType(name, fields) {
    return `${name}(${fields.map(({ name: name2, type }) => type + " " + name2).join(",")})`;
  }
  var TypedDataEncoder = class _TypedDataEncoder {
    /**
     *  The primary type for the structured [[types]].
     *
     *  This is derived automatically from the [[types]], since no
     *  recursion is possible, once the DAG for the types is consturcted
     *  internally, the primary type must be the only remaining type with
     *  no parent nodes.
     */
    primaryType;
    #types;
    /**
     *  The types.
     */
    get types() {
      return JSON.parse(this.#types);
    }
    #fullTypes;
    #encoderCache;
    /**
     *  Create a new **TypedDataEncoder** for %%types%%.
     *
     *  This performs all necessary checking that types are valid and
     *  do not violate the [[link-eip-712]] structural constraints as
     *  well as computes the [[primaryType]].
     */
    constructor(types) {
      this.#types = JSON.stringify(types);
      this.#fullTypes = /* @__PURE__ */ new Map();
      this.#encoderCache = /* @__PURE__ */ new Map();
      const links = /* @__PURE__ */ new Map();
      const parents = /* @__PURE__ */ new Map();
      const subtypes = /* @__PURE__ */ new Map();
      Object.keys(types).forEach((type) => {
        links.set(type, /* @__PURE__ */ new Set());
        parents.set(type, []);
        subtypes.set(type, /* @__PURE__ */ new Set());
      });
      for (const name in types) {
        const uniqueNames = /* @__PURE__ */ new Set();
        for (const field of types[name]) {
          assertArgument(!uniqueNames.has(field.name), `duplicate variable name ${JSON.stringify(field.name)} in ${JSON.stringify(name)}`, "types", types);
          uniqueNames.add(field.name);
          const baseType = field.type.match(/^([^\x5b]*)(\x5b|$)/)[1] || null;
          assertArgument(baseType !== name, `circular type reference to ${JSON.stringify(baseType)}`, "types", types);
          const encoder = getBaseEncoder(baseType);
          if (encoder) {
            continue;
          }
          assertArgument(parents.has(baseType), `unknown type ${JSON.stringify(baseType)}`, "types", types);
          parents.get(baseType).push(name);
          links.get(name).add(baseType);
        }
      }
      const primaryTypes = Array.from(parents.keys()).filter((n4) => parents.get(n4).length === 0);
      assertArgument(primaryTypes.length !== 0, "missing primary type", "types", types);
      assertArgument(primaryTypes.length === 1, `ambiguous primary types or unused types: ${primaryTypes.map((t4) => JSON.stringify(t4)).join(", ")}`, "types", types);
      defineProperties(this, { primaryType: primaryTypes[0] });
      function checkCircular(type, found) {
        assertArgument(!found.has(type), `circular type reference to ${JSON.stringify(type)}`, "types", types);
        found.add(type);
        for (const child of links.get(type)) {
          if (!parents.has(child)) {
            continue;
          }
          checkCircular(child, found);
          for (const subtype of found) {
            subtypes.get(subtype).add(child);
          }
        }
        found.delete(type);
      }
      checkCircular(this.primaryType, /* @__PURE__ */ new Set());
      for (const [name, set] of subtypes) {
        const st = Array.from(set);
        st.sort();
        this.#fullTypes.set(name, encodeType(name, types[name]) + st.map((t4) => encodeType(t4, types[t4])).join(""));
      }
    }
    /**
     *  Returnthe encoder for the specific %%type%%.
     */
    getEncoder(type) {
      let encoder = this.#encoderCache.get(type);
      if (!encoder) {
        encoder = this.#getEncoder(type);
        this.#encoderCache.set(type, encoder);
      }
      return encoder;
    }
    #getEncoder(type) {
      {
        const encoder = getBaseEncoder(type);
        if (encoder) {
          return encoder;
        }
      }
      const match = type.match(/^(.*)(\x5b(\d*)\x5d)$/);
      if (match) {
        const subtype = match[1];
        const subEncoder = this.getEncoder(subtype);
        return (value) => {
          assertArgument(!match[3] || parseInt(match[3]) === value.length, `array length mismatch; expected length ${parseInt(match[3])}`, "value", value);
          let result = value.map(subEncoder);
          if (this.#fullTypes.has(subtype)) {
            result = result.map(keccak256);
          }
          return keccak256(concat(result));
        };
      }
      const fields = this.types[type];
      if (fields) {
        const encodedType = id(this.#fullTypes.get(type));
        return (value) => {
          const values = fields.map(({ name, type: type2 }) => {
            const result = this.getEncoder(type2)(value[name]);
            if (this.#fullTypes.has(type2)) {
              return keccak256(result);
            }
            return result;
          });
          values.unshift(encodedType);
          return concat(values);
        };
      }
      assertArgument(false, `unknown type: ${type}`, "type", type);
    }
    /**
     *  Return the full type for %%name%%.
     */
    encodeType(name) {
      const result = this.#fullTypes.get(name);
      assertArgument(result, `unknown type: ${JSON.stringify(name)}`, "name", name);
      return result;
    }
    /**
     *  Return the encoded %%value%% for the %%type%%.
     */
    encodeData(type, value) {
      return this.getEncoder(type)(value);
    }
    /**
     *  Returns the hash of %%value%% for the type of %%name%%.
     */
    hashStruct(name, value) {
      return keccak256(this.encodeData(name, value));
    }
    /**
     *  Return the fulled encoded %%value%% for the [[types]].
     */
    encode(value) {
      return this.encodeData(this.primaryType, value);
    }
    /**
     *  Return the hash of the fully encoded %%value%% for the [[types]].
     */
    hash(value) {
      return this.hashStruct(this.primaryType, value);
    }
    /**
     *  @_ignore:
     */
    _visit(type, value, callback) {
      {
        const encoder = getBaseEncoder(type);
        if (encoder) {
          return callback(type, value);
        }
      }
      const match = type.match(/^(.*)(\x5b(\d*)\x5d)$/);
      if (match) {
        assertArgument(!match[3] || parseInt(match[3]) === value.length, `array length mismatch; expected length ${parseInt(match[3])}`, "value", value);
        return value.map((v3) => this._visit(match[1], v3, callback));
      }
      const fields = this.types[type];
      if (fields) {
        return fields.reduce((accum, { name, type: type2 }) => {
          accum[name] = this._visit(type2, value[name], callback);
          return accum;
        }, {});
      }
      assertArgument(false, `unknown type: ${type}`, "type", type);
    }
    /**
     *  Call %%calback%% for each value in %%value%%, passing the type and
     *  component within %%value%%.
     *
     *  This is useful for replacing addresses or other transformation that
     *  may be desired on each component, based on its type.
     */
    visit(value, callback) {
      return this._visit(this.primaryType, value, callback);
    }
    /**
     *  Create a new **TypedDataEncoder** for %%types%%.
     */
    static from(types) {
      return new _TypedDataEncoder(types);
    }
    /**
     *  Return the primary type for %%types%%.
     */
    static getPrimaryType(types) {
      return _TypedDataEncoder.from(types).primaryType;
    }
    /**
     *  Return the hashed struct for %%value%% using %%types%% and %%name%%.
     */
    static hashStruct(name, types, value) {
      return _TypedDataEncoder.from(types).hashStruct(name, value);
    }
    /**
     *  Return the domain hash for %%domain%%.
     */
    static hashDomain(domain) {
      const domainFields = [];
      for (const name in domain) {
        if (domain[name] == null) {
          continue;
        }
        const type = domainFieldTypes[name];
        assertArgument(type, `invalid typed-data domain key: ${JSON.stringify(name)}`, "domain", domain);
        domainFields.push({ name, type });
      }
      domainFields.sort((a4, b5) => {
        return domainFieldNames.indexOf(a4.name) - domainFieldNames.indexOf(b5.name);
      });
      return _TypedDataEncoder.hashStruct("EIP712Domain", { EIP712Domain: domainFields }, domain);
    }
    /**
     *  Return the fully encoded [[link-eip-712]] %%value%% for %%types%% with %%domain%%.
     */
    static encode(domain, types, value) {
      return concat([
        "0x1901",
        _TypedDataEncoder.hashDomain(domain),
        _TypedDataEncoder.from(types).hash(value)
      ]);
    }
    /**
     *  Return the hash of the fully encoded [[link-eip-712]] %%value%% for %%types%% with %%domain%%.
     */
    static hash(domain, types, value) {
      return keccak256(_TypedDataEncoder.encode(domain, types, value));
    }
    // Replaces all address types with ENS names with their looked up address
    /**
     * Resolves to the value from resolving all addresses in %%value%% for
     * %%types%% and the %%domain%%.
     */
    static async resolveNames(domain, types, value, resolveName) {
      domain = Object.assign({}, domain);
      for (const key in domain) {
        if (domain[key] == null) {
          delete domain[key];
        }
      }
      const ensCache = {};
      if (domain.verifyingContract && !isHexString(domain.verifyingContract, 20)) {
        ensCache[domain.verifyingContract] = "0x";
      }
      const encoder = _TypedDataEncoder.from(types);
      encoder.visit(value, (type, value2) => {
        if (type === "address" && !isHexString(value2, 20)) {
          ensCache[value2] = "0x";
        }
        return value2;
      });
      for (const name in ensCache) {
        ensCache[name] = await resolveName(name);
      }
      if (domain.verifyingContract && ensCache[domain.verifyingContract]) {
        domain.verifyingContract = ensCache[domain.verifyingContract];
      }
      value = encoder.visit(value, (type, value2) => {
        if (type === "address" && ensCache[value2]) {
          return ensCache[value2];
        }
        return value2;
      });
      return { domain, value };
    }
    /**
     *  Returns the JSON-encoded payload expected by nodes which implement
     *  the JSON-RPC [[link-eip-712]] method.
     */
    static getPayload(domain, types, value) {
      _TypedDataEncoder.hashDomain(domain);
      const domainValues = {};
      const domainTypes = [];
      domainFieldNames.forEach((name) => {
        const value2 = domain[name];
        if (value2 == null) {
          return;
        }
        domainValues[name] = domainChecks[name](value2);
        domainTypes.push({ name, type: domainFieldTypes[name] });
      });
      const encoder = _TypedDataEncoder.from(types);
      const typesWithDomain = Object.assign({}, types);
      assertArgument(typesWithDomain.EIP712Domain == null, "types must not contain EIP712Domain type", "types.EIP712Domain", types);
      typesWithDomain.EIP712Domain = domainTypes;
      encoder.encode(value);
      return {
        types: typesWithDomain,
        domain: domainValues,
        primaryType: encoder.primaryType,
        message: encoder.visit(value, (type, value2) => {
          if (type.match(/^bytes(\d*)/)) {
            return hexlify(getBytes(value2));
          }
          if (type.match(/^u?int/)) {
            return getBigInt(value2).toString();
          }
          switch (type) {
            case "address":
              return value2.toLowerCase();
            case "bool":
              return !!value2;
            case "string":
              assertArgument(typeof value2 === "string", "invalid string", "value", value2);
              return value2;
          }
          assertArgument(false, "unsupported type", "type", type);
        })
      };
    }
  };

  // node_modules/ethers/lib.esm/abi/fragments.js
  function setify(items) {
    const result = /* @__PURE__ */ new Set();
    items.forEach((k3) => result.add(k3));
    return Object.freeze(result);
  }
  var _kwVisibDeploy = "external public payable";
  var KwVisibDeploy = setify(_kwVisibDeploy.split(" "));
  var _kwVisib = "constant external internal payable private public pure view";
  var KwVisib = setify(_kwVisib.split(" "));
  var _kwTypes = "constructor error event fallback function receive struct";
  var KwTypes = setify(_kwTypes.split(" "));
  var _kwModifiers = "calldata memory storage payable indexed";
  var KwModifiers = setify(_kwModifiers.split(" "));
  var _kwOther = "tuple returns";
  var _keywords = [_kwTypes, _kwModifiers, _kwOther, _kwVisib].join(" ");
  var Keywords = setify(_keywords.split(" "));
  var SimpleTokens = {
    "(": "OPEN_PAREN",
    ")": "CLOSE_PAREN",
    "[": "OPEN_BRACKET",
    "]": "CLOSE_BRACKET",
    ",": "COMMA",
    "@": "AT"
  };
  var regexWhitespacePrefix = new RegExp("^(\\s*)");
  var regexNumberPrefix = new RegExp("^([0-9]+)");
  var regexIdPrefix = new RegExp("^([a-zA-Z$_][a-zA-Z0-9$_]*)");
  var regexId = new RegExp("^([a-zA-Z$_][a-zA-Z0-9$_]*)$");
  var regexType = new RegExp("^(address|bool|bytes([0-9]*)|string|u?int([0-9]*))$");
  var TokenString = class _TokenString {
    #offset;
    #tokens;
    get offset() {
      return this.#offset;
    }
    get length() {
      return this.#tokens.length - this.#offset;
    }
    constructor(tokens) {
      this.#offset = 0;
      this.#tokens = tokens.slice();
    }
    clone() {
      return new _TokenString(this.#tokens);
    }
    reset() {
      this.#offset = 0;
    }
    #subTokenString(from = 0, to = 0) {
      return new _TokenString(this.#tokens.slice(from, to).map((t4) => {
        return Object.freeze(Object.assign({}, t4, {
          match: t4.match - from,
          linkBack: t4.linkBack - from,
          linkNext: t4.linkNext - from
        }));
      }));
    }
    // Pops and returns the value of the next token, if it is a keyword in allowed; throws if out of tokens
    popKeyword(allowed) {
      const top = this.peek();
      if (top.type !== "KEYWORD" || !allowed.has(top.text)) {
        throw new Error(`expected keyword ${top.text}`);
      }
      return this.pop().text;
    }
    // Pops and returns the value of the next token if it is `type`; throws if out of tokens
    popType(type) {
      if (this.peek().type !== type) {
        throw new Error(`expected ${type}; got ${JSON.stringify(this.peek())}`);
      }
      return this.pop().text;
    }
    // Pops and returns a "(" TOKENS ")"
    popParen() {
      const top = this.peek();
      if (top.type !== "OPEN_PAREN") {
        throw new Error("bad start");
      }
      const result = this.#subTokenString(this.#offset + 1, top.match + 1);
      this.#offset = top.match + 1;
      return result;
    }
    // Pops and returns the items within "(" ITEM1 "," ITEM2 "," ... ")"
    popParams() {
      const top = this.peek();
      if (top.type !== "OPEN_PAREN") {
        throw new Error("bad start");
      }
      const result = [];
      while (this.#offset < top.match - 1) {
        const link = this.peek().linkNext;
        result.push(this.#subTokenString(this.#offset + 1, link));
        this.#offset = link;
      }
      this.#offset = top.match + 1;
      return result;
    }
    // Returns the top Token, throwing if out of tokens
    peek() {
      if (this.#offset >= this.#tokens.length) {
        throw new Error("out-of-bounds");
      }
      return this.#tokens[this.#offset];
    }
    // Returns the next value, if it is a keyword in `allowed`
    peekKeyword(allowed) {
      const top = this.peekType("KEYWORD");
      return top != null && allowed.has(top) ? top : null;
    }
    // Returns the value of the next token if it is `type`
    peekType(type) {
      if (this.length === 0) {
        return null;
      }
      const top = this.peek();
      return top.type === type ? top.text : null;
    }
    // Returns the next token; throws if out of tokens
    pop() {
      const result = this.peek();
      this.#offset++;
      return result;
    }
    toString() {
      const tokens = [];
      for (let i4 = this.#offset; i4 < this.#tokens.length; i4++) {
        const token = this.#tokens[i4];
        tokens.push(`${token.type}:${token.text}`);
      }
      return `<TokenString ${tokens.join(" ")}>`;
    }
  };
  function lex(text) {
    const tokens = [];
    const throwError2 = (message) => {
      const token = offset < text.length ? JSON.stringify(text[offset]) : "$EOI";
      throw new Error(`invalid token ${token} at ${offset}: ${message}`);
    };
    let brackets = [];
    let commas = [];
    let offset = 0;
    while (offset < text.length) {
      let cur = text.substring(offset);
      let match = cur.match(regexWhitespacePrefix);
      if (match) {
        offset += match[1].length;
        cur = text.substring(offset);
      }
      const token = { depth: brackets.length, linkBack: -1, linkNext: -1, match: -1, type: "", text: "", offset, value: -1 };
      tokens.push(token);
      let type = SimpleTokens[cur[0]] || "";
      if (type) {
        token.type = type;
        token.text = cur[0];
        offset++;
        if (type === "OPEN_PAREN") {
          brackets.push(tokens.length - 1);
          commas.push(tokens.length - 1);
        } else if (type == "CLOSE_PAREN") {
          if (brackets.length === 0) {
            throwError2("no matching open bracket");
          }
          token.match = brackets.pop();
          tokens[token.match].match = tokens.length - 1;
          token.depth--;
          token.linkBack = commas.pop();
          tokens[token.linkBack].linkNext = tokens.length - 1;
        } else if (type === "COMMA") {
          token.linkBack = commas.pop();
          tokens[token.linkBack].linkNext = tokens.length - 1;
          commas.push(tokens.length - 1);
        } else if (type === "OPEN_BRACKET") {
          token.type = "BRACKET";
        } else if (type === "CLOSE_BRACKET") {
          let suffix = tokens.pop().text;
          if (tokens.length > 0 && tokens[tokens.length - 1].type === "NUMBER") {
            const value = tokens.pop().text;
            suffix = value + suffix;
            tokens[tokens.length - 1].value = getNumber(value);
          }
          if (tokens.length === 0 || tokens[tokens.length - 1].type !== "BRACKET") {
            throw new Error("missing opening bracket");
          }
          tokens[tokens.length - 1].text += suffix;
        }
        continue;
      }
      match = cur.match(regexIdPrefix);
      if (match) {
        token.text = match[1];
        offset += token.text.length;
        if (Keywords.has(token.text)) {
          token.type = "KEYWORD";
          continue;
        }
        if (token.text.match(regexType)) {
          token.type = "TYPE";
          continue;
        }
        token.type = "ID";
        continue;
      }
      match = cur.match(regexNumberPrefix);
      if (match) {
        token.text = match[1];
        token.type = "NUMBER";
        offset += token.text.length;
        continue;
      }
      throw new Error(`unexpected token ${JSON.stringify(cur[0])} at position ${offset}`);
    }
    return new TokenString(tokens.map((t4) => Object.freeze(t4)));
  }
  function allowSingle(set, allowed) {
    let included = [];
    for (const key in allowed.keys()) {
      if (set.has(key)) {
        included.push(key);
      }
    }
    if (included.length > 1) {
      throw new Error(`conflicting types: ${included.join(", ")}`);
    }
  }
  function consumeName(type, tokens) {
    if (tokens.peekKeyword(KwTypes)) {
      const keyword = tokens.pop().text;
      if (keyword !== type) {
        throw new Error(`expected ${type}, got ${keyword}`);
      }
    }
    return tokens.popType("ID");
  }
  function consumeKeywords(tokens, allowed) {
    const keywords = /* @__PURE__ */ new Set();
    while (true) {
      const keyword = tokens.peekType("KEYWORD");
      if (keyword == null || allowed && !allowed.has(keyword)) {
        break;
      }
      tokens.pop();
      if (keywords.has(keyword)) {
        throw new Error(`duplicate keywords: ${JSON.stringify(keyword)}`);
      }
      keywords.add(keyword);
    }
    return Object.freeze(keywords);
  }
  function consumeMutability(tokens) {
    let modifiers = consumeKeywords(tokens, KwVisib);
    allowSingle(modifiers, setify("constant payable nonpayable".split(" ")));
    allowSingle(modifiers, setify("pure view payable nonpayable".split(" ")));
    if (modifiers.has("view")) {
      return "view";
    }
    if (modifiers.has("pure")) {
      return "pure";
    }
    if (modifiers.has("payable")) {
      return "payable";
    }
    if (modifiers.has("nonpayable")) {
      return "nonpayable";
    }
    if (modifiers.has("constant")) {
      return "view";
    }
    return "nonpayable";
  }
  function consumeParams(tokens, allowIndexed) {
    return tokens.popParams().map((t4) => ParamType.from(t4, allowIndexed));
  }
  function consumeGas(tokens) {
    if (tokens.peekType("AT")) {
      tokens.pop();
      if (tokens.peekType("NUMBER")) {
        return getBigInt(tokens.pop().text);
      }
      throw new Error("invalid gas");
    }
    return null;
  }
  function consumeEoi(tokens) {
    if (tokens.length) {
      throw new Error(`unexpected tokens: ${tokens.toString()}`);
    }
  }
  var regexArrayType = new RegExp(/^(.*)\[([0-9]*)\]$/);
  function verifyBasicType(type) {
    const match = type.match(regexType);
    assertArgument(match, "invalid type", "type", type);
    if (type === "uint") {
      return "uint256";
    }
    if (type === "int") {
      return "int256";
    }
    if (match[2]) {
      const length = parseInt(match[2]);
      assertArgument(length !== 0 && length <= 32, "invalid bytes length", "type", type);
    } else if (match[3]) {
      const size = parseInt(match[3]);
      assertArgument(size !== 0 && size <= 256 && size % 8 === 0, "invalid numeric width", "type", type);
    }
    return type;
  }
  var _guard4 = {};
  var internal = Symbol.for("_ethers_internal");
  var ParamTypeInternal = "_ParamTypeInternal";
  var ErrorFragmentInternal = "_ErrorInternal";
  var EventFragmentInternal = "_EventInternal";
  var ConstructorFragmentInternal = "_ConstructorInternal";
  var FallbackFragmentInternal = "_FallbackInternal";
  var FunctionFragmentInternal = "_FunctionInternal";
  var StructFragmentInternal = "_StructInternal";
  var ParamType = class _ParamType {
    /**
     *  The local name of the parameter (or ``""`` if unbound)
     */
    name;
    /**
     *  The fully qualified type (e.g. ``"address"``, ``"tuple(address)"``,
     *  ``"uint256[3][]"``)
     */
    type;
    /**
     *  The base type (e.g. ``"address"``, ``"tuple"``, ``"array"``)
     */
    baseType;
    /**
     *  True if the parameters is indexed.
     *
     *  For non-indexable types this is ``null``.
     */
    indexed;
    /**
     *  The components for the tuple.
     *
     *  For non-tuple types this is ``null``.
     */
    components;
    /**
     *  The array length, or ``-1`` for dynamic-lengthed arrays.
     *
     *  For non-array types this is ``null``.
     */
    arrayLength;
    /**
     *  The type of each child in the array.
     *
     *  For non-array types this is ``null``.
     */
    arrayChildren;
    /**
     *  @private
     */
    constructor(guard, name, type, baseType, indexed, components, arrayLength, arrayChildren) {
      assertPrivate(guard, _guard4, "ParamType");
      Object.defineProperty(this, internal, { value: ParamTypeInternal });
      if (components) {
        components = Object.freeze(components.slice());
      }
      if (baseType === "array") {
        if (arrayLength == null || arrayChildren == null) {
          throw new Error("");
        }
      } else if (arrayLength != null || arrayChildren != null) {
        throw new Error("");
      }
      if (baseType === "tuple") {
        if (components == null) {
          throw new Error("");
        }
      } else if (components != null) {
        throw new Error("");
      }
      defineProperties(this, {
        name,
        type,
        baseType,
        indexed,
        components,
        arrayLength,
        arrayChildren
      });
    }
    /**
     *  Return a string representation of this type.
     *
     *  For example,
     *
     *  ``sighash" => "(uint256,address)"``
     *
     *  ``"minimal" => "tuple(uint256,address) indexed"``
     *
     *  ``"full" => "tuple(uint256 foo, address bar) indexed baz"``
     */
    format(format) {
      if (format == null) {
        format = "sighash";
      }
      if (format === "json") {
        const name = this.name || "";
        if (this.isArray()) {
          const result3 = JSON.parse(this.arrayChildren.format("json"));
          result3.name = name;
          result3.type += `[${this.arrayLength < 0 ? "" : String(this.arrayLength)}]`;
          return JSON.stringify(result3);
        }
        const result2 = {
          type: this.baseType === "tuple" ? "tuple" : this.type,
          name
        };
        if (typeof this.indexed === "boolean") {
          result2.indexed = this.indexed;
        }
        if (this.isTuple()) {
          result2.components = this.components.map((c4) => JSON.parse(c4.format(format)));
        }
        return JSON.stringify(result2);
      }
      let result = "";
      if (this.isArray()) {
        result += this.arrayChildren.format(format);
        result += `[${this.arrayLength < 0 ? "" : String(this.arrayLength)}]`;
      } else {
        if (this.isTuple()) {
          if (format !== "sighash") {
            result += this.type;
          }
          result += "(" + this.components.map((comp) => comp.format(format)).join(format === "full" ? ", " : ",") + ")";
        } else {
          result += this.type;
        }
      }
      if (format !== "sighash") {
        if (this.indexed === true) {
          result += " indexed";
        }
        if (format === "full" && this.name) {
          result += " " + this.name;
        }
      }
      return result;
    }
    /**
     *  Returns true if %%this%% is an Array type.
     *
     *  This provides a type gaurd ensuring that [[arrayChildren]]
     *  and [[arrayLength]] are non-null.
     */
    isArray() {
      return this.baseType === "array";
    }
    /**
     *  Returns true if %%this%% is a Tuple type.
     *
     *  This provides a type gaurd ensuring that [[components]]
     *  is non-null.
     */
    isTuple() {
      return this.baseType === "tuple";
    }
    /**
     *  Returns true if %%this%% is an Indexable type.
     *
     *  This provides a type gaurd ensuring that [[indexed]]
     *  is non-null.
     */
    isIndexable() {
      return this.indexed != null;
    }
    /**
     *  Walks the **ParamType** with %%value%%, calling %%process%%
     *  on each type, destructing the %%value%% recursively.
     */
    walk(value, process2) {
      if (this.isArray()) {
        if (!Array.isArray(value)) {
          throw new Error("invalid array value");
        }
        if (this.arrayLength !== -1 && value.length !== this.arrayLength) {
          throw new Error("array is wrong length");
        }
        const _this = this;
        return value.map((v3) => _this.arrayChildren.walk(v3, process2));
      }
      if (this.isTuple()) {
        if (!Array.isArray(value)) {
          throw new Error("invalid tuple value");
        }
        if (value.length !== this.components.length) {
          throw new Error("array is wrong length");
        }
        const _this = this;
        return value.map((v3, i4) => _this.components[i4].walk(v3, process2));
      }
      return process2(this.type, value);
    }
    #walkAsync(promises, value, process2, setValue) {
      if (this.isArray()) {
        if (!Array.isArray(value)) {
          throw new Error("invalid array value");
        }
        if (this.arrayLength !== -1 && value.length !== this.arrayLength) {
          throw new Error("array is wrong length");
        }
        const childType = this.arrayChildren;
        const result2 = value.slice();
        result2.forEach((value2, index) => {
          childType.#walkAsync(promises, value2, process2, (value3) => {
            result2[index] = value3;
          });
        });
        setValue(result2);
        return;
      }
      if (this.isTuple()) {
        const components = this.components;
        let result2;
        if (Array.isArray(value)) {
          result2 = value.slice();
        } else {
          if (value == null || typeof value !== "object") {
            throw new Error("invalid tuple value");
          }
          result2 = components.map((param) => {
            if (!param.name) {
              throw new Error("cannot use object value with unnamed components");
            }
            if (!(param.name in value)) {
              throw new Error(`missing value for component ${param.name}`);
            }
            return value[param.name];
          });
        }
        if (result2.length !== this.components.length) {
          throw new Error("array is wrong length");
        }
        result2.forEach((value2, index) => {
          components[index].#walkAsync(promises, value2, process2, (value3) => {
            result2[index] = value3;
          });
        });
        setValue(result2);
        return;
      }
      const result = process2(this.type, value);
      if (result.then) {
        promises.push(async function() {
          setValue(await result);
        }());
      } else {
        setValue(result);
      }
    }
    /**
     *  Walks the **ParamType** with %%value%%, asynchronously calling
     *  %%process%% on each type, destructing the %%value%% recursively.
     *
     *  This can be used to resolve ENS naes by walking and resolving each
     *  ``"address"`` type.
     */
    async walkAsync(value, process2) {
      const promises = [];
      const result = [value];
      this.#walkAsync(promises, value, process2, (value2) => {
        result[0] = value2;
      });
      if (promises.length) {
        await Promise.all(promises);
      }
      return result[0];
    }
    /**
     *  Creates a new **ParamType** for %%obj%%.
     *
     *  If %%allowIndexed%% then the ``indexed`` keyword is permitted,
     *  otherwise the ``indexed`` keyword will throw an error.
     */
    static from(obj, allowIndexed) {
      if (_ParamType.isParamType(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        try {
          return _ParamType.from(lex(obj), allowIndexed);
        } catch (error) {
          assertArgument(false, "invalid param type", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        let type2 = "", baseType = "";
        let comps = null;
        if (consumeKeywords(obj, setify(["tuple"])).has("tuple") || obj.peekType("OPEN_PAREN")) {
          baseType = "tuple";
          comps = obj.popParams().map((t4) => _ParamType.from(t4));
          type2 = `tuple(${comps.map((c4) => c4.format()).join(",")})`;
        } else {
          type2 = verifyBasicType(obj.popType("TYPE"));
          baseType = type2;
        }
        let arrayChildren = null;
        let arrayLength = null;
        while (obj.length && obj.peekType("BRACKET")) {
          const bracket = obj.pop();
          arrayChildren = new _ParamType(_guard4, "", type2, baseType, null, comps, arrayLength, arrayChildren);
          arrayLength = bracket.value;
          type2 += bracket.text;
          baseType = "array";
          comps = null;
        }
        let indexed2 = null;
        const keywords = consumeKeywords(obj, KwModifiers);
        if (keywords.has("indexed")) {
          if (!allowIndexed) {
            throw new Error("");
          }
          indexed2 = true;
        }
        const name2 = obj.peekType("ID") ? obj.pop().text : "";
        if (obj.length) {
          throw new Error("leftover tokens");
        }
        return new _ParamType(_guard4, name2, type2, baseType, indexed2, comps, arrayLength, arrayChildren);
      }
      const name = obj.name;
      assertArgument(!name || typeof name === "string" && name.match(regexId), "invalid name", "obj.name", name);
      let indexed = obj.indexed;
      if (indexed != null) {
        assertArgument(allowIndexed, "parameter cannot be indexed", "obj.indexed", obj.indexed);
        indexed = !!indexed;
      }
      let type = obj.type;
      let arrayMatch = type.match(regexArrayType);
      if (arrayMatch) {
        const arrayLength = parseInt(arrayMatch[2] || "-1");
        const arrayChildren = _ParamType.from({
          type: arrayMatch[1],
          components: obj.components
        });
        return new _ParamType(_guard4, name || "", type, "array", indexed, null, arrayLength, arrayChildren);
      }
      if (type === "tuple" || type.startsWith(
        "tuple("
        /* fix: ) */
      ) || type.startsWith(
        "("
        /* fix: ) */
      )) {
        const comps = obj.components != null ? obj.components.map((c4) => _ParamType.from(c4)) : null;
        const tuple = new _ParamType(_guard4, name || "", type, "tuple", indexed, comps, null, null);
        return tuple;
      }
      type = verifyBasicType(obj.type);
      return new _ParamType(_guard4, name || "", type, type, indexed, null, null, null);
    }
    /**
     *  Returns true if %%value%% is a **ParamType**.
     */
    static isParamType(value) {
      return value && value[internal] === ParamTypeInternal;
    }
  };
  var Fragment = class _Fragment {
    /**
     *  The type of the fragment.
     */
    type;
    /**
     *  The inputs for the fragment.
     */
    inputs;
    /**
     *  @private
     */
    constructor(guard, type, inputs) {
      assertPrivate(guard, _guard4, "Fragment");
      inputs = Object.freeze(inputs.slice());
      defineProperties(this, { type, inputs });
    }
    /**
     *  Creates a new **Fragment** for %%obj%%, wich can be any supported
     *  ABI frgament type.
     */
    static from(obj) {
      if (typeof obj === "string") {
        try {
          _Fragment.from(JSON.parse(obj));
        } catch (e4) {
        }
        return _Fragment.from(lex(obj));
      }
      if (obj instanceof TokenString) {
        const type = obj.peekKeyword(KwTypes);
        switch (type) {
          case "constructor":
            return ConstructorFragment.from(obj);
          case "error":
            return ErrorFragment.from(obj);
          case "event":
            return EventFragment.from(obj);
          case "fallback":
          case "receive":
            return FallbackFragment.from(obj);
          case "function":
            return FunctionFragment.from(obj);
          case "struct":
            return StructFragment.from(obj);
        }
      } else if (typeof obj === "object") {
        switch (obj.type) {
          case "constructor":
            return ConstructorFragment.from(obj);
          case "error":
            return ErrorFragment.from(obj);
          case "event":
            return EventFragment.from(obj);
          case "fallback":
          case "receive":
            return FallbackFragment.from(obj);
          case "function":
            return FunctionFragment.from(obj);
          case "struct":
            return StructFragment.from(obj);
        }
        assert(false, `unsupported type: ${obj.type}`, "UNSUPPORTED_OPERATION", {
          operation: "Fragment.from"
        });
      }
      assertArgument(false, "unsupported frgament object", "obj", obj);
    }
    /**
     *  Returns true if %%value%% is a [[ConstructorFragment]].
     */
    static isConstructor(value) {
      return ConstructorFragment.isFragment(value);
    }
    /**
     *  Returns true if %%value%% is an [[ErrorFragment]].
     */
    static isError(value) {
      return ErrorFragment.isFragment(value);
    }
    /**
     *  Returns true if %%value%% is an [[EventFragment]].
     */
    static isEvent(value) {
      return EventFragment.isFragment(value);
    }
    /**
     *  Returns true if %%value%% is a [[FunctionFragment]].
     */
    static isFunction(value) {
      return FunctionFragment.isFragment(value);
    }
    /**
     *  Returns true if %%value%% is a [[StructFragment]].
     */
    static isStruct(value) {
      return StructFragment.isFragment(value);
    }
  };
  var NamedFragment = class extends Fragment {
    /**
     *  The name of the fragment.
     */
    name;
    /**
     *  @private
     */
    constructor(guard, type, name, inputs) {
      super(guard, type, inputs);
      assertArgument(typeof name === "string" && name.match(regexId), "invalid identifier", "name", name);
      inputs = Object.freeze(inputs.slice());
      defineProperties(this, { name });
    }
  };
  function joinParams(format, params) {
    return "(" + params.map((p4) => p4.format(format)).join(format === "full" ? ", " : ",") + ")";
  }
  var ErrorFragment = class _ErrorFragment extends NamedFragment {
    /**
     *  @private
     */
    constructor(guard, name, inputs) {
      super(guard, "error", name, inputs);
      Object.defineProperty(this, internal, { value: ErrorFragmentInternal });
    }
    /**
     *  The Custom Error selector.
     */
    get selector() {
      return id(this.format("sighash")).substring(0, 10);
    }
    /**
     *  Returns a string representation of this fragment as %%format%%.
     */
    format(format) {
      if (format == null) {
        format = "sighash";
      }
      if (format === "json") {
        return JSON.stringify({
          type: "error",
          name: this.name,
          inputs: this.inputs.map((input) => JSON.parse(input.format(format)))
        });
      }
      const result = [];
      if (format !== "sighash") {
        result.push("error");
      }
      result.push(this.name + joinParams(format, this.inputs));
      return result.join(" ");
    }
    /**
     *  Returns a new **ErrorFragment** for %%obj%%.
     */
    static from(obj) {
      if (_ErrorFragment.isFragment(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        return _ErrorFragment.from(lex(obj));
      } else if (obj instanceof TokenString) {
        const name = consumeName("error", obj);
        const inputs = consumeParams(obj);
        consumeEoi(obj);
        return new _ErrorFragment(_guard4, name, inputs);
      }
      return new _ErrorFragment(_guard4, obj.name, obj.inputs ? obj.inputs.map(ParamType.from) : []);
    }
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is an
     *  **ErrorFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === ErrorFragmentInternal;
    }
  };
  var EventFragment = class _EventFragment extends NamedFragment {
    /**
     *  Whether this event is anonymous.
     */
    anonymous;
    /**
     *  @private
     */
    constructor(guard, name, inputs, anonymous) {
      super(guard, "event", name, inputs);
      Object.defineProperty(this, internal, { value: EventFragmentInternal });
      defineProperties(this, { anonymous });
    }
    /**
     *  The Event topic hash.
     */
    get topicHash() {
      return id(this.format("sighash"));
    }
    /**
     *  Returns a string representation of this event as %%format%%.
     */
    format(format) {
      if (format == null) {
        format = "sighash";
      }
      if (format === "json") {
        return JSON.stringify({
          type: "event",
          anonymous: this.anonymous,
          name: this.name,
          inputs: this.inputs.map((i4) => JSON.parse(i4.format(format)))
        });
      }
      const result = [];
      if (format !== "sighash") {
        result.push("event");
      }
      result.push(this.name + joinParams(format, this.inputs));
      if (format !== "sighash" && this.anonymous) {
        result.push("anonymous");
      }
      return result.join(" ");
    }
    /**
     *  Return the topic hash for an event with %%name%% and %%params%%.
     */
    static getTopicHash(name, params) {
      params = (params || []).map((p4) => ParamType.from(p4));
      const fragment = new _EventFragment(_guard4, name, params, false);
      return fragment.topicHash;
    }
    /**
     *  Returns a new **EventFragment** for %%obj%%.
     */
    static from(obj) {
      if (_EventFragment.isFragment(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        try {
          return _EventFragment.from(lex(obj));
        } catch (error) {
          assertArgument(false, "invalid event fragment", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        const name = consumeName("event", obj);
        const inputs = consumeParams(obj, true);
        const anonymous = !!consumeKeywords(obj, setify(["anonymous"])).has("anonymous");
        consumeEoi(obj);
        return new _EventFragment(_guard4, name, inputs, anonymous);
      }
      return new _EventFragment(_guard4, obj.name, obj.inputs ? obj.inputs.map((p4) => ParamType.from(p4, true)) : [], !!obj.anonymous);
    }
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is an
     *  **EventFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === EventFragmentInternal;
    }
  };
  var ConstructorFragment = class _ConstructorFragment extends Fragment {
    /**
     *  Whether the constructor can receive an endowment.
     */
    payable;
    /**
     *  The recommended gas limit for deployment or ``null``.
     */
    gas;
    /**
     *  @private
     */
    constructor(guard, type, inputs, payable, gas) {
      super(guard, type, inputs);
      Object.defineProperty(this, internal, { value: ConstructorFragmentInternal });
      defineProperties(this, { payable, gas });
    }
    /**
     *  Returns a string representation of this constructor as %%format%%.
     */
    format(format) {
      assert(format != null && format !== "sighash", "cannot format a constructor for sighash", "UNSUPPORTED_OPERATION", { operation: "format(sighash)" });
      if (format === "json") {
        return JSON.stringify({
          type: "constructor",
          stateMutability: this.payable ? "payable" : "undefined",
          payable: this.payable,
          gas: this.gas != null ? this.gas : void 0,
          inputs: this.inputs.map((i4) => JSON.parse(i4.format(format)))
        });
      }
      const result = [`constructor${joinParams(format, this.inputs)}`];
      result.push(this.payable ? "payable" : "nonpayable");
      if (this.gas != null) {
        result.push(`@${this.gas.toString()}`);
      }
      return result.join(" ");
    }
    /**
     *  Returns a new **ConstructorFragment** for %%obj%%.
     */
    static from(obj) {
      if (_ConstructorFragment.isFragment(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        try {
          return _ConstructorFragment.from(lex(obj));
        } catch (error) {
          assertArgument(false, "invalid constuctor fragment", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        consumeKeywords(obj, setify(["constructor"]));
        const inputs = consumeParams(obj);
        const payable = !!consumeKeywords(obj, KwVisibDeploy).has("payable");
        const gas = consumeGas(obj);
        consumeEoi(obj);
        return new _ConstructorFragment(_guard4, "constructor", inputs, payable, gas);
      }
      return new _ConstructorFragment(_guard4, "constructor", obj.inputs ? obj.inputs.map(ParamType.from) : [], !!obj.payable, obj.gas != null ? obj.gas : null);
    }
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is a
     *  **ConstructorFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === ConstructorFragmentInternal;
    }
  };
  var FallbackFragment = class _FallbackFragment extends Fragment {
    /**
     *  If the function can be sent value during invocation.
     */
    payable;
    constructor(guard, inputs, payable) {
      super(guard, "fallback", inputs);
      Object.defineProperty(this, internal, { value: FallbackFragmentInternal });
      defineProperties(this, { payable });
    }
    /**
     *  Returns a string representation of this fallback as %%format%%.
     */
    format(format) {
      const type = this.inputs.length === 0 ? "receive" : "fallback";
      if (format === "json") {
        const stateMutability = this.payable ? "payable" : "nonpayable";
        return JSON.stringify({ type, stateMutability });
      }
      return `${type}()${this.payable ? " payable" : ""}`;
    }
    /**
     *  Returns a new **FallbackFragment** for %%obj%%.
     */
    static from(obj) {
      if (_FallbackFragment.isFragment(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        try {
          return _FallbackFragment.from(lex(obj));
        } catch (error) {
          assertArgument(false, "invalid fallback fragment", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        const errorObj = obj.toString();
        const topIsValid = obj.peekKeyword(setify(["fallback", "receive"]));
        assertArgument(topIsValid, "type must be fallback or receive", "obj", errorObj);
        const type = obj.popKeyword(setify(["fallback", "receive"]));
        if (type === "receive") {
          const inputs2 = consumeParams(obj);
          assertArgument(inputs2.length === 0, `receive cannot have arguments`, "obj.inputs", inputs2);
          consumeKeywords(obj, setify(["payable"]));
          consumeEoi(obj);
          return new _FallbackFragment(_guard4, [], true);
        }
        let inputs = consumeParams(obj);
        if (inputs.length) {
          assertArgument(inputs.length === 1 && inputs[0].type === "bytes", "invalid fallback inputs", "obj.inputs", inputs.map((i4) => i4.format("minimal")).join(", "));
        } else {
          inputs = [ParamType.from("bytes")];
        }
        const mutability = consumeMutability(obj);
        assertArgument(mutability === "nonpayable" || mutability === "payable", "fallback cannot be constants", "obj.stateMutability", mutability);
        if (consumeKeywords(obj, setify(["returns"])).has("returns")) {
          const outputs = consumeParams(obj);
          assertArgument(outputs.length === 1 && outputs[0].type === "bytes", "invalid fallback outputs", "obj.outputs", outputs.map((i4) => i4.format("minimal")).join(", "));
        }
        consumeEoi(obj);
        return new _FallbackFragment(_guard4, inputs, mutability === "payable");
      }
      if (obj.type === "receive") {
        return new _FallbackFragment(_guard4, [], true);
      }
      if (obj.type === "fallback") {
        const inputs = [ParamType.from("bytes")];
        const payable = obj.stateMutability === "payable";
        return new _FallbackFragment(_guard4, inputs, payable);
      }
      assertArgument(false, "invalid fallback description", "obj", obj);
    }
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is a
     *  **FallbackFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === FallbackFragmentInternal;
    }
  };
  var FunctionFragment = class _FunctionFragment extends NamedFragment {
    /**
     *  If the function is constant (e.g. ``pure`` or ``view`` functions).
     */
    constant;
    /**
     *  The returned types for the result of calling this function.
     */
    outputs;
    /**
     *  The state mutability (e.g. ``payable``, ``nonpayable``, ``view``
     *  or ``pure``)
     */
    stateMutability;
    /**
     *  If the function can be sent value during invocation.
     */
    payable;
    /**
     *  The recommended gas limit to send when calling this function.
     */
    gas;
    /**
     *  @private
     */
    constructor(guard, name, stateMutability, inputs, outputs, gas) {
      super(guard, "function", name, inputs);
      Object.defineProperty(this, internal, { value: FunctionFragmentInternal });
      outputs = Object.freeze(outputs.slice());
      const constant = stateMutability === "view" || stateMutability === "pure";
      const payable = stateMutability === "payable";
      defineProperties(this, { constant, gas, outputs, payable, stateMutability });
    }
    /**
     *  The Function selector.
     */
    get selector() {
      return id(this.format("sighash")).substring(0, 10);
    }
    /**
     *  Returns a string representation of this function as %%format%%.
     */
    format(format) {
      if (format == null) {
        format = "sighash";
      }
      if (format === "json") {
        return JSON.stringify({
          type: "function",
          name: this.name,
          constant: this.constant,
          stateMutability: this.stateMutability !== "nonpayable" ? this.stateMutability : void 0,
          payable: this.payable,
          gas: this.gas != null ? this.gas : void 0,
          inputs: this.inputs.map((i4) => JSON.parse(i4.format(format))),
          outputs: this.outputs.map((o4) => JSON.parse(o4.format(format)))
        });
      }
      const result = [];
      if (format !== "sighash") {
        result.push("function");
      }
      result.push(this.name + joinParams(format, this.inputs));
      if (format !== "sighash") {
        if (this.stateMutability !== "nonpayable") {
          result.push(this.stateMutability);
        }
        if (this.outputs && this.outputs.length) {
          result.push("returns");
          result.push(joinParams(format, this.outputs));
        }
        if (this.gas != null) {
          result.push(`@${this.gas.toString()}`);
        }
      }
      return result.join(" ");
    }
    /**
     *  Return the selector for a function with %%name%% and %%params%%.
     */
    static getSelector(name, params) {
      params = (params || []).map((p4) => ParamType.from(p4));
      const fragment = new _FunctionFragment(_guard4, name, "view", params, [], null);
      return fragment.selector;
    }
    /**
     *  Returns a new **FunctionFragment** for %%obj%%.
     */
    static from(obj) {
      if (_FunctionFragment.isFragment(obj)) {
        return obj;
      }
      if (typeof obj === "string") {
        try {
          return _FunctionFragment.from(lex(obj));
        } catch (error) {
          assertArgument(false, "invalid function fragment", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        const name = consumeName("function", obj);
        const inputs = consumeParams(obj);
        const mutability = consumeMutability(obj);
        let outputs = [];
        if (consumeKeywords(obj, setify(["returns"])).has("returns")) {
          outputs = consumeParams(obj);
        }
        const gas = consumeGas(obj);
        consumeEoi(obj);
        return new _FunctionFragment(_guard4, name, mutability, inputs, outputs, gas);
      }
      let stateMutability = obj.stateMutability;
      if (stateMutability == null) {
        stateMutability = "payable";
        if (typeof obj.constant === "boolean") {
          stateMutability = "view";
          if (!obj.constant) {
            stateMutability = "payable";
            if (typeof obj.payable === "boolean" && !obj.payable) {
              stateMutability = "nonpayable";
            }
          }
        } else if (typeof obj.payable === "boolean" && !obj.payable) {
          stateMutability = "nonpayable";
        }
      }
      return new _FunctionFragment(_guard4, obj.name, stateMutability, obj.inputs ? obj.inputs.map(ParamType.from) : [], obj.outputs ? obj.outputs.map(ParamType.from) : [], obj.gas != null ? obj.gas : null);
    }
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is a
     *  **FunctionFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === FunctionFragmentInternal;
    }
  };
  var StructFragment = class _StructFragment extends NamedFragment {
    /**
     *  @private
     */
    constructor(guard, name, inputs) {
      super(guard, "struct", name, inputs);
      Object.defineProperty(this, internal, { value: StructFragmentInternal });
    }
    /**
     *  Returns a string representation of this struct as %%format%%.
     */
    format() {
      throw new Error("@TODO");
    }
    /**
     *  Returns a new **StructFragment** for %%obj%%.
     */
    static from(obj) {
      if (typeof obj === "string") {
        try {
          return _StructFragment.from(lex(obj));
        } catch (error) {
          assertArgument(false, "invalid struct fragment", "obj", obj);
        }
      } else if (obj instanceof TokenString) {
        const name = consumeName("struct", obj);
        const inputs = consumeParams(obj);
        consumeEoi(obj);
        return new _StructFragment(_guard4, name, inputs);
      }
      return new _StructFragment(_guard4, obj.name, obj.inputs ? obj.inputs.map(ParamType.from) : []);
    }
    // @TODO: fix this return type
    /**
     *  Returns ``true`` and provides a type guard if %%value%% is a
     *  **StructFragment**.
     */
    static isFragment(value) {
      return value && value[internal] === StructFragmentInternal;
    }
  };

  // node_modules/ethers/lib.esm/abi/abi-coder.js
  var PanicReasons = /* @__PURE__ */ new Map();
  PanicReasons.set(0, "GENERIC_PANIC");
  PanicReasons.set(1, "ASSERT_FALSE");
  PanicReasons.set(17, "OVERFLOW");
  PanicReasons.set(18, "DIVIDE_BY_ZERO");
  PanicReasons.set(33, "ENUM_RANGE_ERROR");
  PanicReasons.set(34, "BAD_STORAGE_DATA");
  PanicReasons.set(49, "STACK_UNDERFLOW");
  PanicReasons.set(50, "ARRAY_RANGE_ERROR");
  PanicReasons.set(65, "OUT_OF_MEMORY");
  PanicReasons.set(81, "UNINITIALIZED_FUNCTION_CALL");
  var paramTypeBytes = new RegExp(/^bytes([0-9]*)$/);
  var paramTypeNumber = new RegExp(/^(u?int)([0-9]*)$/);
  var defaultCoder = null;
  function getBuiltinCallException(action, tx, data, abiCoder) {
    let message = "missing revert data";
    let reason = null;
    const invocation = null;
    let revert = null;
    if (data) {
      message = "execution reverted";
      const bytes2 = getBytes(data);
      data = hexlify(data);
      if (bytes2.length === 0) {
        message += " (no data present; likely require(false) occurred";
        reason = "require(false)";
      } else if (bytes2.length % 32 !== 4) {
        message += " (could not decode reason; invalid data length)";
      } else if (hexlify(bytes2.slice(0, 4)) === "0x08c379a0") {
        try {
          reason = abiCoder.decode(["string"], bytes2.slice(4))[0];
          revert = {
            signature: "Error(string)",
            name: "Error",
            args: [reason]
          };
          message += `: ${JSON.stringify(reason)}`;
        } catch (error) {
          message += " (could not decode reason; invalid string data)";
        }
      } else if (hexlify(bytes2.slice(0, 4)) === "0x4e487b71") {
        try {
          const code = Number(abiCoder.decode(["uint256"], bytes2.slice(4))[0]);
          revert = {
            signature: "Panic(uint256)",
            name: "Panic",
            args: [code]
          };
          reason = `Panic due to ${PanicReasons.get(code) || "UNKNOWN"}(${code})`;
          message += `: ${reason}`;
        } catch (error) {
          message += " (could not decode panic code)";
        }
      } else {
        message += " (unknown custom error)";
      }
    }
    const transaction = {
      to: tx.to ? getAddress(tx.to) : null,
      data: tx.data || "0x"
    };
    if (tx.from) {
      transaction.from = getAddress(tx.from);
    }
    return makeError(message, "CALL_EXCEPTION", {
      action,
      data,
      reason,
      transaction,
      invocation,
      revert
    });
  }
  var AbiCoder = class _AbiCoder {
    #getCoder(param) {
      if (param.isArray()) {
        return new ArrayCoder(this.#getCoder(param.arrayChildren), param.arrayLength, param.name);
      }
      if (param.isTuple()) {
        return new TupleCoder(param.components.map((c4) => this.#getCoder(c4)), param.name);
      }
      switch (param.baseType) {
        case "address":
          return new AddressCoder(param.name);
        case "bool":
          return new BooleanCoder(param.name);
        case "string":
          return new StringCoder(param.name);
        case "bytes":
          return new BytesCoder(param.name);
        case "":
          return new NullCoder(param.name);
      }
      let match = param.type.match(paramTypeNumber);
      if (match) {
        let size = parseInt(match[2] || "256");
        assertArgument(size !== 0 && size <= 256 && size % 8 === 0, "invalid " + match[1] + " bit length", "param", param);
        return new NumberCoder(size / 8, match[1] === "int", param.name);
      }
      match = param.type.match(paramTypeBytes);
      if (match) {
        let size = parseInt(match[1]);
        assertArgument(size !== 0 && size <= 32, "invalid bytes length", "param", param);
        return new FixedBytesCoder(size, param.name);
      }
      assertArgument(false, "invalid type", "type", param.type);
    }
    /**
     *  Get the default values for the given %%types%%.
     *
     *  For example, a ``uint`` is by default ``0`` and ``bool``
     *  is by default ``false``.
     */
    getDefaultValue(types) {
      const coders = types.map((type) => this.#getCoder(ParamType.from(type)));
      const coder = new TupleCoder(coders, "_");
      return coder.defaultValue();
    }
    /**
     *  Encode the %%values%% as the %%types%% into ABI data.
     *
     *  @returns DataHexstring
     */
    encode(types, values) {
      assertArgumentCount(values.length, types.length, "types/values length mismatch");
      const coders = types.map((type) => this.#getCoder(ParamType.from(type)));
      const coder = new TupleCoder(coders, "_");
      const writer = new Writer();
      coder.encode(writer, values);
      return writer.data;
    }
    /**
     *  Decode the ABI %%data%% as the %%types%% into values.
     *
     *  If %%loose%% decoding is enabled, then strict padding is
     *  not enforced. Some older versions of Solidity incorrectly
     *  padded event data emitted from ``external`` functions.
     */
    decode(types, data, loose) {
      const coders = types.map((type) => this.#getCoder(ParamType.from(type)));
      const coder = new TupleCoder(coders, "_");
      return coder.decode(new Reader(data, loose));
    }
    /**
     *  Returns the shared singleton instance of a default [[AbiCoder]].
     *
     *  On the first call, the instance is created internally.
     */
    static defaultAbiCoder() {
      if (defaultCoder == null) {
        defaultCoder = new _AbiCoder();
      }
      return defaultCoder;
    }
    /**
     *  Returns an ethers-compatible [[CallExceptionError]] Error for the given
     *  result %%data%% for the [[CallExceptionAction]] %%action%% against
     *  the Transaction %%tx%%.
     */
    static getBuiltinCallException(action, tx, data) {
      return getBuiltinCallException(action, tx, data, _AbiCoder.defaultAbiCoder());
    }
  };

  // node_modules/ethers/lib.esm/abi/interface.js
  var LogDescription = class {
    /**
     *  The matching fragment for the ``topic0``.
     */
    fragment;
    /**
     *  The name of the Event.
     */
    name;
    /**
     *  The full Event signature.
     */
    signature;
    /**
     *  The topic hash for the Event.
     */
    topic;
    /**
     *  The arguments passed into the Event with ``emit``.
     */
    args;
    /**
     *  @_ignore:
     */
    constructor(fragment, topic, args) {
      const name = fragment.name, signature = fragment.format();
      defineProperties(this, {
        fragment,
        name,
        signature,
        topic,
        args
      });
    }
  };
  var TransactionDescription = class {
    /**
     *  The matching fragment from the transaction ``data``.
     */
    fragment;
    /**
     *  The name of the Function from the transaction ``data``.
     */
    name;
    /**
     *  The arguments passed to the Function from the transaction ``data``.
     */
    args;
    /**
     *  The full Function signature from the transaction ``data``.
     */
    signature;
    /**
     *  The selector for the Function from the transaction ``data``.
     */
    selector;
    /**
     *  The ``value`` (in wei) from the transaction.
     */
    value;
    /**
     *  @_ignore:
     */
    constructor(fragment, selector, args, value) {
      const name = fragment.name, signature = fragment.format();
      defineProperties(this, {
        fragment,
        name,
        args,
        signature,
        selector,
        value
      });
    }
  };
  var ErrorDescription = class {
    /**
     *  The matching fragment.
     */
    fragment;
    /**
     *  The name of the Error.
     */
    name;
    /**
     *  The arguments passed to the Error with ``revert``.
     */
    args;
    /**
     *  The full Error signature.
     */
    signature;
    /**
     *  The selector for the Error.
     */
    selector;
    /**
     *  @_ignore:
     */
    constructor(fragment, selector, args) {
      const name = fragment.name, signature = fragment.format();
      defineProperties(this, {
        fragment,
        name,
        args,
        signature,
        selector
      });
    }
  };
  var Indexed = class {
    /**
     *  The ``keccak256`` of the value logged.
     */
    hash;
    /**
     *  @_ignore:
     */
    _isIndexed;
    /**
     *  Returns ``true`` if %%value%% is an **Indexed**.
     *
     *  This provides a Type Guard for property access.
     */
    static isIndexed(value) {
      return !!(value && value._isIndexed);
    }
    /**
     *  @_ignore:
     */
    constructor(hash2) {
      defineProperties(this, { hash: hash2, _isIndexed: true });
    }
  };
  var PanicReasons2 = {
    "0": "generic panic",
    "1": "assert(false)",
    "17": "arithmetic overflow",
    "18": "division or modulo by zero",
    "33": "enum overflow",
    "34": "invalid encoded storage byte array accessed",
    "49": "out-of-bounds array access; popping on an empty array",
    "50": "out-of-bounds access of an array or bytesN",
    "65": "out of memory",
    "81": "uninitialized function"
  };
  var BuiltinErrors = {
    "0x08c379a0": {
      signature: "Error(string)",
      name: "Error",
      inputs: ["string"],
      reason: (message) => {
        return `reverted with reason string ${JSON.stringify(message)}`;
      }
    },
    "0x4e487b71": {
      signature: "Panic(uint256)",
      name: "Panic",
      inputs: ["uint256"],
      reason: (code) => {
        let reason = "unknown panic code";
        if (code >= 0 && code <= 255 && PanicReasons2[code.toString()]) {
          reason = PanicReasons2[code.toString()];
        }
        return `reverted with panic code 0x${code.toString(16)} (${reason})`;
      }
    }
  };
  var Interface = class _Interface {
    /**
     *  All the Contract ABI members (i.e. methods, events, errors, etc).
     */
    fragments;
    /**
     *  The Contract constructor.
     */
    deploy;
    /**
     *  The Fallback method, if any.
     */
    fallback;
    /**
     *  If receiving ether is supported.
     */
    receive;
    #errors;
    #events;
    #functions;
    //    #structs: Map<string, StructFragment>;
    #abiCoder;
    /**
     *  Create a new Interface for the %%fragments%%.
     */
    constructor(fragments) {
      let abi = [];
      if (typeof fragments === "string") {
        abi = JSON.parse(fragments);
      } else {
        abi = fragments;
      }
      this.#functions = /* @__PURE__ */ new Map();
      this.#errors = /* @__PURE__ */ new Map();
      this.#events = /* @__PURE__ */ new Map();
      const frags = [];
      for (const a4 of abi) {
        try {
          frags.push(Fragment.from(a4));
        } catch (error) {
          console.log("EE", error);
        }
      }
      defineProperties(this, {
        fragments: Object.freeze(frags)
      });
      let fallback = null;
      let receive = false;
      this.#abiCoder = this.getAbiCoder();
      this.fragments.forEach((fragment, index) => {
        let bucket;
        switch (fragment.type) {
          case "constructor":
            if (this.deploy) {
              console.log("duplicate definition - constructor");
              return;
            }
            defineProperties(this, { deploy: fragment });
            return;
          case "fallback":
            if (fragment.inputs.length === 0) {
              receive = true;
            } else {
              assertArgument(!fallback || fragment.payable !== fallback.payable, "conflicting fallback fragments", `fragments[${index}]`, fragment);
              fallback = fragment;
              receive = fallback.payable;
            }
            return;
          case "function":
            bucket = this.#functions;
            break;
          case "event":
            bucket = this.#events;
            break;
          case "error":
            bucket = this.#errors;
            break;
          default:
            return;
        }
        const signature = fragment.format();
        if (bucket.has(signature)) {
          return;
        }
        bucket.set(signature, fragment);
      });
      if (!this.deploy) {
        defineProperties(this, {
          deploy: ConstructorFragment.from("constructor()")
        });
      }
      defineProperties(this, { fallback, receive });
    }
    /**
     *  Returns the entire Human-Readable ABI, as an array of
     *  signatures, optionally as %%minimal%% strings, which
     *  removes parameter names and unneceesary spaces.
     */
    format(minimal) {
      const format = minimal ? "minimal" : "full";
      const abi = this.fragments.map((f4) => f4.format(format));
      return abi;
    }
    /**
     *  Return the JSON-encoded ABI. This is the format Solidiy
     *  returns.
     */
    formatJson() {
      const abi = this.fragments.map((f4) => f4.format("json"));
      return JSON.stringify(abi.map((j4) => JSON.parse(j4)));
    }
    /**
     *  The ABI coder that will be used to encode and decode binary
     *  data.
     */
    getAbiCoder() {
      return AbiCoder.defaultAbiCoder();
    }
    // Find a function definition by any means necessary (unless it is ambiguous)
    #getFunction(key, values, forceUnique) {
      if (isHexString(key)) {
        const selector = key.toLowerCase();
        for (const fragment of this.#functions.values()) {
          if (selector === fragment.selector) {
            return fragment;
          }
        }
        return null;
      }
      if (key.indexOf("(") === -1) {
        const matching = [];
        for (const [name, fragment] of this.#functions) {
          if (name.split(
            "("
            /* fix:) */
          )[0] === key) {
            matching.push(fragment);
          }
        }
        if (values) {
          const lastValue = values.length > 0 ? values[values.length - 1] : null;
          let valueLength = values.length;
          let allowOptions = true;
          if (Typed.isTyped(lastValue) && lastValue.type === "overrides") {
            allowOptions = false;
            valueLength--;
          }
          for (let i4 = matching.length - 1; i4 >= 0; i4--) {
            const inputs = matching[i4].inputs.length;
            if (inputs !== valueLength && (!allowOptions || inputs !== valueLength - 1)) {
              matching.splice(i4, 1);
            }
          }
          for (let i4 = matching.length - 1; i4 >= 0; i4--) {
            const inputs = matching[i4].inputs;
            for (let j4 = 0; j4 < values.length; j4++) {
              if (!Typed.isTyped(values[j4])) {
                continue;
              }
              if (j4 >= inputs.length) {
                if (values[j4].type === "overrides") {
                  continue;
                }
                matching.splice(i4, 1);
                break;
              }
              if (values[j4].type !== inputs[j4].baseType) {
                matching.splice(i4, 1);
                break;
              }
            }
          }
        }
        if (matching.length === 1 && values && values.length !== matching[0].inputs.length) {
          const lastArg = values[values.length - 1];
          if (lastArg == null || Array.isArray(lastArg) || typeof lastArg !== "object") {
            matching.splice(0, 1);
          }
        }
        if (matching.length === 0) {
          return null;
        }
        if (matching.length > 1 && forceUnique) {
          const matchStr = matching.map((m4) => JSON.stringify(m4.format())).join(", ");
          assertArgument(false, `ambiguous function description (i.e. matches ${matchStr})`, "key", key);
        }
        return matching[0];
      }
      const result = this.#functions.get(FunctionFragment.from(key).format());
      if (result) {
        return result;
      }
      return null;
    }
    /**
     *  Get the function name for %%key%%, which may be a function selector,
     *  function name or function signature that belongs to the ABI.
     */
    getFunctionName(key) {
      const fragment = this.#getFunction(key, null, false);
      assertArgument(fragment, "no matching function", "key", key);
      return fragment.name;
    }
    /**
     *  Returns true if %%key%% (a function selector, function name or
     *  function signature) is present in the ABI.
     *
     *  In the case of a function name, the name may be ambiguous, so
     *  accessing the [[FunctionFragment]] may require refinement.
     */
    hasFunction(key) {
      return !!this.#getFunction(key, null, false);
    }
    /**
     *  Get the [[FunctionFragment]] for %%key%%, which may be a function
     *  selector, function name or function signature that belongs to the ABI.
     *
     *  If %%values%% is provided, it will use the Typed API to handle
     *  ambiguous cases where multiple functions match by name.
     *
     *  If the %%key%% and %%values%% do not refine to a single function in
     *  the ABI, this will throw.
     */
    getFunction(key, values) {
      return this.#getFunction(key, values || null, true);
    }
    /**
     *  Iterate over all functions, calling %%callback%%, sorted by their name.
     */
    forEachFunction(callback) {
      const names2 = Array.from(this.#functions.keys());
      names2.sort((a4, b5) => a4.localeCompare(b5));
      for (let i4 = 0; i4 < names2.length; i4++) {
        const name = names2[i4];
        callback(this.#functions.get(name), i4);
      }
    }
    // Find an event definition by any means necessary (unless it is ambiguous)
    #getEvent(key, values, forceUnique) {
      if (isHexString(key)) {
        const eventTopic = key.toLowerCase();
        for (const fragment of this.#events.values()) {
          if (eventTopic === fragment.topicHash) {
            return fragment;
          }
        }
        return null;
      }
      if (key.indexOf("(") === -1) {
        const matching = [];
        for (const [name, fragment] of this.#events) {
          if (name.split(
            "("
            /* fix:) */
          )[0] === key) {
            matching.push(fragment);
          }
        }
        if (values) {
          for (let i4 = matching.length - 1; i4 >= 0; i4--) {
            if (matching[i4].inputs.length < values.length) {
              matching.splice(i4, 1);
            }
          }
          for (let i4 = matching.length - 1; i4 >= 0; i4--) {
            const inputs = matching[i4].inputs;
            for (let j4 = 0; j4 < values.length; j4++) {
              if (!Typed.isTyped(values[j4])) {
                continue;
              }
              if (values[j4].type !== inputs[j4].baseType) {
                matching.splice(i4, 1);
                break;
              }
            }
          }
        }
        if (matching.length === 0) {
          return null;
        }
        if (matching.length > 1 && forceUnique) {
          const matchStr = matching.map((m4) => JSON.stringify(m4.format())).join(", ");
          assertArgument(false, `ambiguous event description (i.e. matches ${matchStr})`, "key", key);
        }
        return matching[0];
      }
      const result = this.#events.get(EventFragment.from(key).format());
      if (result) {
        return result;
      }
      return null;
    }
    /**
     *  Get the event name for %%key%%, which may be a topic hash,
     *  event name or event signature that belongs to the ABI.
     */
    getEventName(key) {
      const fragment = this.#getEvent(key, null, false);
      assertArgument(fragment, "no matching event", "key", key);
      return fragment.name;
    }
    /**
     *  Returns true if %%key%% (an event topic hash, event name or
     *  event signature) is present in the ABI.
     *
     *  In the case of an event name, the name may be ambiguous, so
     *  accessing the [[EventFragment]] may require refinement.
     */
    hasEvent(key) {
      return !!this.#getEvent(key, null, false);
    }
    /**
     *  Get the [[EventFragment]] for %%key%%, which may be a topic hash,
     *  event name or event signature that belongs to the ABI.
     *
     *  If %%values%% is provided, it will use the Typed API to handle
     *  ambiguous cases where multiple events match by name.
     *
     *  If the %%key%% and %%values%% do not refine to a single event in
     *  the ABI, this will throw.
     */
    getEvent(key, values) {
      return this.#getEvent(key, values || null, true);
    }
    /**
     *  Iterate over all events, calling %%callback%%, sorted by their name.
     */
    forEachEvent(callback) {
      const names2 = Array.from(this.#events.keys());
      names2.sort((a4, b5) => a4.localeCompare(b5));
      for (let i4 = 0; i4 < names2.length; i4++) {
        const name = names2[i4];
        callback(this.#events.get(name), i4);
      }
    }
    /**
     *  Get the [[ErrorFragment]] for %%key%%, which may be an error
     *  selector, error name or error signature that belongs to the ABI.
     *
     *  If %%values%% is provided, it will use the Typed API to handle
     *  ambiguous cases where multiple errors match by name.
     *
     *  If the %%key%% and %%values%% do not refine to a single error in
     *  the ABI, this will throw.
     */
    getError(key, values) {
      if (isHexString(key)) {
        const selector = key.toLowerCase();
        if (BuiltinErrors[selector]) {
          return ErrorFragment.from(BuiltinErrors[selector].signature);
        }
        for (const fragment of this.#errors.values()) {
          if (selector === fragment.selector) {
            return fragment;
          }
        }
        return null;
      }
      if (key.indexOf("(") === -1) {
        const matching = [];
        for (const [name, fragment] of this.#errors) {
          if (name.split(
            "("
            /* fix:) */
          )[0] === key) {
            matching.push(fragment);
          }
        }
        if (matching.length === 0) {
          if (key === "Error") {
            return ErrorFragment.from("error Error(string)");
          }
          if (key === "Panic") {
            return ErrorFragment.from("error Panic(uint256)");
          }
          return null;
        } else if (matching.length > 1) {
          const matchStr = matching.map((m4) => JSON.stringify(m4.format())).join(", ");
          assertArgument(false, `ambiguous error description (i.e. ${matchStr})`, "name", key);
        }
        return matching[0];
      }
      key = ErrorFragment.from(key).format();
      if (key === "Error(string)") {
        return ErrorFragment.from("error Error(string)");
      }
      if (key === "Panic(uint256)") {
        return ErrorFragment.from("error Panic(uint256)");
      }
      const result = this.#errors.get(key);
      if (result) {
        return result;
      }
      return null;
    }
    /**
     *  Iterate over all errors, calling %%callback%%, sorted by their name.
     */
    forEachError(callback) {
      const names2 = Array.from(this.#errors.keys());
      names2.sort((a4, b5) => a4.localeCompare(b5));
      for (let i4 = 0; i4 < names2.length; i4++) {
        const name = names2[i4];
        callback(this.#errors.get(name), i4);
      }
    }
    // Get the 4-byte selector used by Solidity to identify a function
    /*
    getSelector(fragment: ErrorFragment | FunctionFragment): string {
        if (typeof(fragment) === "string") {
            const matches: Array<Fragment> = [ ];
    
            try { matches.push(this.getFunction(fragment)); } catch (error) { }
            try { matches.push(this.getError(<string>fragment)); } catch (_) { }
    
            if (matches.length === 0) {
                logger.throwArgumentError("unknown fragment", "key", fragment);
            } else if (matches.length > 1) {
                logger.throwArgumentError("ambiguous fragment matches function and error", "key", fragment);
            }
    
            fragment = matches[0];
        }
    
        return dataSlice(id(fragment.format()), 0, 4);
    }
        */
    // Get the 32-byte topic hash used by Solidity to identify an event
    /*
    getEventTopic(fragment: EventFragment): string {
        //if (typeof(fragment) === "string") { fragment = this.getEvent(eventFragment); }
        return id(fragment.format());
    }
    */
    _decodeParams(params, data) {
      return this.#abiCoder.decode(params, data);
    }
    _encodeParams(params, values) {
      return this.#abiCoder.encode(params, values);
    }
    /**
     *  Encodes a ``tx.data`` object for deploying the Contract with
     *  the %%values%% as the constructor arguments.
     */
    encodeDeploy(values) {
      return this._encodeParams(this.deploy.inputs, values || []);
    }
    /**
     *  Decodes the result %%data%% (e.g. from an ``eth_call``) for the
     *  specified error (see [[getError]] for valid values for
     *  %%key%%).
     *
     *  Most developers should prefer the [[parseCallResult]] method instead,
     *  which will automatically detect a ``CALL_EXCEPTION`` and throw the
     *  corresponding error.
     */
    decodeErrorResult(fragment, data) {
      if (typeof fragment === "string") {
        const f4 = this.getError(fragment);
        assertArgument(f4, "unknown error", "fragment", fragment);
        fragment = f4;
      }
      assertArgument(dataSlice(data, 0, 4) === fragment.selector, `data signature does not match error ${fragment.name}.`, "data", data);
      return this._decodeParams(fragment.inputs, dataSlice(data, 4));
    }
    /**
     *  Encodes the transaction revert data for a call result that
     *  reverted from the the Contract with the sepcified %%error%%
     *  (see [[getError]] for valid values for %%fragment%%) with the %%values%%.
     *
     *  This is generally not used by most developers, unless trying to mock
     *  a result from a Contract.
     */
    encodeErrorResult(fragment, values) {
      if (typeof fragment === "string") {
        const f4 = this.getError(fragment);
        assertArgument(f4, "unknown error", "fragment", fragment);
        fragment = f4;
      }
      return concat([
        fragment.selector,
        this._encodeParams(fragment.inputs, values || [])
      ]);
    }
    /**
     *  Decodes the %%data%% from a transaction ``tx.data`` for
     *  the function specified (see [[getFunction]] for valid values
     *  for %%fragment%%).
     *
     *  Most developers should prefer the [[parseTransaction]] method
     *  instead, which will automatically detect the fragment.
     */
    decodeFunctionData(fragment, data) {
      if (typeof fragment === "string") {
        const f4 = this.getFunction(fragment);
        assertArgument(f4, "unknown function", "fragment", fragment);
        fragment = f4;
      }
      assertArgument(dataSlice(data, 0, 4) === fragment.selector, `data signature does not match function ${fragment.name}.`, "data", data);
      return this._decodeParams(fragment.inputs, dataSlice(data, 4));
    }
    /**
     *  Encodes the ``tx.data`` for a transaction that calls the function
     *  specified (see [[getFunction]] for valid values for %%fragment%%) with
     *  the %%values%%.
     */
    encodeFunctionData(fragment, values) {
      if (typeof fragment === "string") {
        const f4 = this.getFunction(fragment);
        assertArgument(f4, "unknown function", "fragment", fragment);
        fragment = f4;
      }
      return concat([
        fragment.selector,
        this._encodeParams(fragment.inputs, values || [])
      ]);
    }
    /**
     *  Decodes the result %%data%% (e.g. from an ``eth_call``) for the
     *  specified function (see [[getFunction]] for valid values for
     *  %%key%%).
     *
     *  Most developers should prefer the [[parseCallResult]] method instead,
     *  which will automatically detect a ``CALL_EXCEPTION`` and throw the
     *  corresponding error.
     */
    decodeFunctionResult(fragment, data) {
      if (typeof fragment === "string") {
        const f4 = this.getFunction(fragment);
        assertArgument(f4, "unknown function", "fragment", fragment);
        fragment = f4;
      }
      let message = "invalid length for result data";
      const bytes2 = getBytesCopy(data);
      if (bytes2.length % 32 === 0) {
        try {
          return this.#abiCoder.decode(fragment.outputs, bytes2);
        } catch (error) {
          message = "could not decode result data";
        }
      }
      assert(false, message, "BAD_DATA", {
        value: hexlify(bytes2),
        info: { method: fragment.name, signature: fragment.format() }
      });
    }
    makeError(_data, tx) {
      const data = getBytes(_data, "data");
      const error = AbiCoder.getBuiltinCallException("call", tx, data);
      const customPrefix = "execution reverted (unknown custom error)";
      if (error.message.startsWith(customPrefix)) {
        const selector = hexlify(data.slice(0, 4));
        const ef = this.getError(selector);
        if (ef) {
          try {
            const args = this.#abiCoder.decode(ef.inputs, data.slice(4));
            error.revert = {
              name: ef.name,
              signature: ef.format(),
              args
            };
            error.reason = error.revert.signature;
            error.message = `execution reverted: ${error.reason}`;
          } catch (e4) {
            error.message = `execution reverted (coult not decode custom error)`;
          }
        }
      }
      const parsed = this.parseTransaction(tx);
      if (parsed) {
        error.invocation = {
          method: parsed.name,
          signature: parsed.signature,
          args: parsed.args
        };
      }
      return error;
    }
    /**
     *  Encodes the result data (e.g. from an ``eth_call``) for the
     *  specified function (see [[getFunction]] for valid values
     *  for %%fragment%%) with %%values%%.
     *
     *  This is generally not used by most developers, unless trying to mock
     *  a result from a Contract.
     */
    encodeFunctionResult(fragment, values) {
      if (typeof fragment === "string") {
        const f4 = this.getFunction(fragment);
        assertArgument(f4, "unknown function", "fragment", fragment);
        fragment = f4;
      }
      return hexlify(this.#abiCoder.encode(fragment.outputs, values || []));
    }
    /*
        spelunk(inputs: Array<ParamType>, values: ReadonlyArray<any>, processfunc: (type: string, value: any) => Promise<any>): Promise<Array<any>> {
            const promises: Array<Promise<>> = [ ];
            const process = function(type: ParamType, value: any): any {
                if (type.baseType === "array") {
                    return descend(type.child
                }
                if (type. === "address") {
                }
            };
    
            const descend = function (inputs: Array<ParamType>, values: ReadonlyArray<any>) {
                if (inputs.length !== values.length) { throw new Error("length mismatch"); }
                
            };
    
            const result: Array<any> = [ ];
            values.forEach((value, index) => {
                if (value == null) {
                    topics.push(null);
                } else if (param.baseType === "array" || param.baseType === "tuple") {
                    logger.throwArgumentError("filtering with tuples or arrays not supported", ("contract." + param.name), value);
                } else if (Array.isArray(value)) {
                    topics.push(value.map((value) => encodeTopic(param, value)));
                } else {
                    topics.push(encodeTopic(param, value));
                }
            });
        }
    */
    // Create the filter for the event with search criteria (e.g. for eth_filterLog)
    encodeFilterTopics(fragment, values) {
      if (typeof fragment === "string") {
        const f4 = this.getEvent(fragment);
        assertArgument(f4, "unknown event", "eventFragment", fragment);
        fragment = f4;
      }
      assert(values.length <= fragment.inputs.length, `too many arguments for ${fragment.format()}`, "UNEXPECTED_ARGUMENT", { count: values.length, expectedCount: fragment.inputs.length });
      const topics = [];
      if (!fragment.anonymous) {
        topics.push(fragment.topicHash);
      }
      const encodeTopic = (param, value) => {
        if (param.type === "string") {
          return id(value);
        } else if (param.type === "bytes") {
          return keccak256(hexlify(value));
        }
        if (param.type === "bool" && typeof value === "boolean") {
          value = value ? "0x01" : "0x00";
        } else if (param.type.match(/^u?int/)) {
          value = toBeHex(value);
        } else if (param.type.match(/^bytes/)) {
          value = zeroPadBytes(value, 32);
        } else if (param.type === "address") {
          this.#abiCoder.encode(["address"], [value]);
        }
        return zeroPadValue(hexlify(value), 32);
      };
      values.forEach((value, index) => {
        const param = fragment.inputs[index];
        if (!param.indexed) {
          assertArgument(value == null, "cannot filter non-indexed parameters; must be null", "contract." + param.name, value);
          return;
        }
        if (value == null) {
          topics.push(null);
        } else if (param.baseType === "array" || param.baseType === "tuple") {
          assertArgument(false, "filtering with tuples or arrays not supported", "contract." + param.name, value);
        } else if (Array.isArray(value)) {
          topics.push(value.map((value2) => encodeTopic(param, value2)));
        } else {
          topics.push(encodeTopic(param, value));
        }
      });
      while (topics.length && topics[topics.length - 1] === null) {
        topics.pop();
      }
      return topics;
    }
    encodeEventLog(fragment, values) {
      if (typeof fragment === "string") {
        const f4 = this.getEvent(fragment);
        assertArgument(f4, "unknown event", "eventFragment", fragment);
        fragment = f4;
      }
      const topics = [];
      const dataTypes = [];
      const dataValues = [];
      if (!fragment.anonymous) {
        topics.push(fragment.topicHash);
      }
      assertArgument(values.length === fragment.inputs.length, "event arguments/values mismatch", "values", values);
      fragment.inputs.forEach((param, index) => {
        const value = values[index];
        if (param.indexed) {
          if (param.type === "string") {
            topics.push(id(value));
          } else if (param.type === "bytes") {
            topics.push(keccak256(value));
          } else if (param.baseType === "tuple" || param.baseType === "array") {
            throw new Error("not implemented");
          } else {
            topics.push(this.#abiCoder.encode([param.type], [value]));
          }
        } else {
          dataTypes.push(param);
          dataValues.push(value);
        }
      });
      return {
        data: this.#abiCoder.encode(dataTypes, dataValues),
        topics
      };
    }
    // Decode a filter for the event and the search criteria
    decodeEventLog(fragment, data, topics) {
      if (typeof fragment === "string") {
        const f4 = this.getEvent(fragment);
        assertArgument(f4, "unknown event", "eventFragment", fragment);
        fragment = f4;
      }
      if (topics != null && !fragment.anonymous) {
        const eventTopic = fragment.topicHash;
        assertArgument(isHexString(topics[0], 32) && topics[0].toLowerCase() === eventTopic, "fragment/topic mismatch", "topics[0]", topics[0]);
        topics = topics.slice(1);
      }
      const indexed = [];
      const nonIndexed = [];
      const dynamic = [];
      fragment.inputs.forEach((param, index) => {
        if (param.indexed) {
          if (param.type === "string" || param.type === "bytes" || param.baseType === "tuple" || param.baseType === "array") {
            indexed.push(ParamType.from({ type: "bytes32", name: param.name }));
            dynamic.push(true);
          } else {
            indexed.push(param);
            dynamic.push(false);
          }
        } else {
          nonIndexed.push(param);
          dynamic.push(false);
        }
      });
      const resultIndexed = topics != null ? this.#abiCoder.decode(indexed, concat(topics)) : null;
      const resultNonIndexed = this.#abiCoder.decode(nonIndexed, data, true);
      const values = [];
      const keys = [];
      let nonIndexedIndex = 0, indexedIndex = 0;
      fragment.inputs.forEach((param, index) => {
        let value = null;
        if (param.indexed) {
          if (resultIndexed == null) {
            value = new Indexed(null);
          } else if (dynamic[index]) {
            value = new Indexed(resultIndexed[indexedIndex++]);
          } else {
            try {
              value = resultIndexed[indexedIndex++];
            } catch (error) {
              value = error;
            }
          }
        } else {
          try {
            value = resultNonIndexed[nonIndexedIndex++];
          } catch (error) {
            value = error;
          }
        }
        values.push(value);
        keys.push(param.name || null);
      });
      return Result.fromItems(values, keys);
    }
    /**
     *  Parses a transaction, finding the matching function and extracts
     *  the parameter values along with other useful function details.
     *
     *  If the matching function cannot be found, return null.
     */
    parseTransaction(tx) {
      const data = getBytes(tx.data, "tx.data");
      const value = getBigInt(tx.value != null ? tx.value : 0, "tx.value");
      const fragment = this.getFunction(hexlify(data.slice(0, 4)));
      if (!fragment) {
        return null;
      }
      const args = this.#abiCoder.decode(fragment.inputs, data.slice(4));
      return new TransactionDescription(fragment, fragment.selector, args, value);
    }
    parseCallResult(data) {
      throw new Error("@TODO");
    }
    /**
     *  Parses a receipt log, finding the matching event and extracts
     *  the parameter values along with other useful event details.
     *
     *  If the matching event cannot be found, returns null.
     */
    parseLog(log) {
      const fragment = this.getEvent(log.topics[0]);
      if (!fragment || fragment.anonymous) {
        return null;
      }
      return new LogDescription(fragment, fragment.topicHash, this.decodeEventLog(fragment, log.data, log.topics));
    }
    /**
     *  Parses a revert data, finding the matching error and extracts
     *  the parameter values along with other useful error details.
     *
     *  If the matching event cannot be found, returns null.
     */
    parseError(data) {
      const hexData = hexlify(data);
      const fragment = this.getError(dataSlice(hexData, 0, 4));
      if (!fragment) {
        return null;
      }
      const args = this.#abiCoder.decode(fragment.inputs, dataSlice(hexData, 4));
      return new ErrorDescription(fragment, fragment.selector, args);
    }
    /**
     *  Creates a new [[Interface]] from the ABI %%value%%.
     *
     *  The %%value%% may be provided as an existing [[Interface]] object,
     *  a JSON-encoded ABI or any Human-Readable ABI format.
     */
    static from(value) {
      if (value instanceof _Interface) {
        return value;
      }
      if (typeof value === "string") {
        return new _Interface(JSON.parse(value));
      }
      if (typeof value.format === "function") {
        return new _Interface(value.format("json"));
      }
      return new _Interface(value);
    }
  };

  // node_modules/ethers/lib.esm/providers/provider.js
  var BN_09 = BigInt(0);
  function getValue2(value) {
    if (value == null) {
      return null;
    }
    return value;
  }
  function toJson(value) {
    if (value == null) {
      return null;
    }
    return value.toString();
  }
  var FeeData = class {
    /**
     *  The gas price for legacy networks.
     */
    gasPrice;
    /**
     *  The maximum fee to pay per gas.
     *
     *  The base fee per gas is defined by the network and based on
     *  congestion, increasing the cost during times of heavy load
     *  and lowering when less busy.
     *
     *  The actual fee per gas will be the base fee for the block
     *  and the priority fee, up to the max fee per gas.
     *
     *  This will be ``null`` on legacy networks (i.e. [pre-EIP-1559](link-eip-1559))
     */
    maxFeePerGas;
    /**
     *  The additional amout to pay per gas to encourage a validator
     *  to include the transaction.
     *
     *  The purpose of this is to compensate the validator for the
     *  adjusted risk for including a given transaction.
     *
     *  This will be ``null`` on legacy networks (i.e. [pre-EIP-1559](link-eip-1559))
     */
    maxPriorityFeePerGas;
    /**
     *  Creates a new FeeData for %%gasPrice%%, %%maxFeePerGas%% and
     *  %%maxPriorityFeePerGas%%.
     */
    constructor(gasPrice, maxFeePerGas, maxPriorityFeePerGas) {
      defineProperties(this, {
        gasPrice: getValue2(gasPrice),
        maxFeePerGas: getValue2(maxFeePerGas),
        maxPriorityFeePerGas: getValue2(maxPriorityFeePerGas)
      });
    }
    /**
     *  Returns a JSON-friendly value.
     */
    toJSON() {
      const { gasPrice, maxFeePerGas, maxPriorityFeePerGas } = this;
      return {
        _type: "FeeData",
        gasPrice: toJson(gasPrice),
        maxFeePerGas: toJson(maxFeePerGas),
        maxPriorityFeePerGas: toJson(maxPriorityFeePerGas)
      };
    }
  };
  function copyRequest(req) {
    const result = {};
    if (req.to) {
      result.to = req.to;
    }
    if (req.from) {
      result.from = req.from;
    }
    if (req.data) {
      result.data = hexlify(req.data);
    }
    const bigIntKeys = "chainId,gasLimit,gasPrice,maxFeePerGas,maxPriorityFeePerGas,value".split(/,/);
    for (const key of bigIntKeys) {
      if (!(key in req) || req[key] == null) {
        continue;
      }
      result[key] = getBigInt(req[key], `request.${key}`);
    }
    const numberKeys = "type,nonce".split(/,/);
    for (const key of numberKeys) {
      if (!(key in req) || req[key] == null) {
        continue;
      }
      result[key] = getNumber(req[key], `request.${key}`);
    }
    if (req.accessList) {
      result.accessList = accessListify(req.accessList);
    }
    if ("blockTag" in req) {
      result.blockTag = req.blockTag;
    }
    if ("enableCcipRead" in req) {
      result.enableCcipRead = !!req.enableCcipRead;
    }
    if ("customData" in req) {
      result.customData = req.customData;
    }
    return result;
  }
  var Block = class {
    /**
     *  The provider connected to the block used to fetch additional details
     *  if necessary.
     */
    provider;
    /**
     *  The block number, sometimes called the block height. This is a
     *  sequential number that is one higher than the parent block.
     */
    number;
    /**
     *  The block hash.
     *
     *  This hash includes all properties, so can be safely used to identify
     *  an exact set of block properties.
     */
    hash;
    /**
     *  The timestamp for this block, which is the number of seconds since
     *  epoch that this block was included.
     */
    timestamp;
    /**
     *  The block hash of the parent block.
     */
    parentHash;
    /**
     *  The nonce.
     *
     *  On legacy networks, this is the random number inserted which
     *  permitted the difficulty target to be reached.
     */
    nonce;
    /**
     *  The difficulty target.
     *
     *  On legacy networks, this is the proof-of-work target required
     *  for a block to meet the protocol rules to be included.
     *
     *  On modern networks, this is a random number arrived at using
     *  randao.  @TODO: Find links?
     */
    difficulty;
    /**
     *  The total gas limit for this block.
     */
    gasLimit;
    /**
     *  The total gas used in this block.
     */
    gasUsed;
    /**
     *  The miner coinbase address, wihch receives any subsidies for
     *  including this block.
     */
    miner;
    /**
     *  Any extra data the validator wished to include.
     */
    extraData;
    /**
     *  The base fee per gas that all transactions in this block were
     *  charged.
     *
     *  This adjusts after each block, depending on how congested the network
     *  is.
     */
    baseFeePerGas;
    #transactions;
    /**
     *  Create a new **Block** object.
     *
     *  This should generally not be necessary as the unless implementing a
     *  low-level library.
     */
    constructor(block, provider) {
      this.#transactions = block.transactions.map((tx) => {
        if (typeof tx !== "string") {
          return new TransactionResponse(tx, provider);
        }
        return tx;
      });
      defineProperties(this, {
        provider,
        hash: getValue2(block.hash),
        number: block.number,
        timestamp: block.timestamp,
        parentHash: block.parentHash,
        nonce: block.nonce,
        difficulty: block.difficulty,
        gasLimit: block.gasLimit,
        gasUsed: block.gasUsed,
        miner: block.miner,
        extraData: block.extraData,
        baseFeePerGas: getValue2(block.baseFeePerGas)
      });
    }
    /**
     *  Returns the list of transaction hashes.
     */
    get transactions() {
      return this.#transactions.map((tx) => {
        if (typeof tx === "string") {
          return tx;
        }
        return tx.hash;
      });
    }
    /**
     *  Returns the complete transactions for blocks which
     *  prefetched them, by passing ``true`` to %%prefetchTxs%%
     *  into [[Provider-getBlock]].
     */
    get prefetchedTransactions() {
      const txs = this.#transactions.slice();
      if (txs.length === 0) {
        return [];
      }
      assert(typeof txs[0] === "object", "transactions were not prefetched with block request", "UNSUPPORTED_OPERATION", {
        operation: "transactionResponses()"
      });
      return txs;
    }
    /**
     *  Returns a JSON-friendly value.
     */
    toJSON() {
      const { baseFeePerGas, difficulty, extraData, gasLimit, gasUsed, hash: hash2, miner, nonce, number: number2, parentHash, timestamp, transactions } = this;
      return {
        _type: "Block",
        baseFeePerGas: toJson(baseFeePerGas),
        difficulty: toJson(difficulty),
        extraData,
        gasLimit: toJson(gasLimit),
        gasUsed: toJson(gasUsed),
        hash: hash2,
        miner,
        nonce,
        number: number2,
        parentHash,
        timestamp,
        transactions
      };
    }
    [Symbol.iterator]() {
      let index = 0;
      const txs = this.transactions;
      return {
        next: () => {
          if (index < this.length) {
            return {
              value: txs[index++],
              done: false
            };
          }
          return { value: void 0, done: true };
        }
      };
    }
    /**
     *  The number of transactions in this block.
     */
    get length() {
      return this.#transactions.length;
    }
    /**
     *  The [[link-js-date]] this block was included at.
     */
    get date() {
      if (this.timestamp == null) {
        return null;
      }
      return new Date(this.timestamp * 1e3);
    }
    /**
     *  Get the transaction at %%indexe%% within this block.
     */
    async getTransaction(indexOrHash) {
      let tx = void 0;
      if (typeof indexOrHash === "number") {
        tx = this.#transactions[indexOrHash];
      } else {
        const hash2 = indexOrHash.toLowerCase();
        for (const v3 of this.#transactions) {
          if (typeof v3 === "string") {
            if (v3 !== hash2) {
              continue;
            }
            tx = v3;
            break;
          } else {
            if (v3.hash === hash2) {
              continue;
            }
            tx = v3;
            break;
          }
        }
      }
      if (tx == null) {
        throw new Error("no such tx");
      }
      if (typeof tx === "string") {
        return await this.provider.getTransaction(tx);
      } else {
        return tx;
      }
    }
    /**
     *  If a **Block** was fetched with a request to include the transactions
     *  this will allow synchronous access to those transactions.
     *
     *  If the transactions were not prefetched, this will throw.
     */
    getPrefetchedTransaction(indexOrHash) {
      const txs = this.prefetchedTransactions;
      if (typeof indexOrHash === "number") {
        return txs[indexOrHash];
      }
      indexOrHash = indexOrHash.toLowerCase();
      for (const tx of txs) {
        if (tx.hash === indexOrHash) {
          return tx;
        }
      }
      assertArgument(false, "no matching transaction", "indexOrHash", indexOrHash);
    }
    /**
     *  Returns true if this block been mined. This provides a type guard
     *  for all properties on a [[MinedBlock]].
     */
    isMined() {
      return !!this.hash;
    }
    /**
     *  Returns true if this block is an [[link-eip-2930]] block.
     */
    isLondon() {
      return !!this.baseFeePerGas;
    }
    /**
     *  @_ignore:
     */
    orphanedEvent() {
      if (!this.isMined()) {
        throw new Error("");
      }
      return createOrphanedBlockFilter(this);
    }
  };
  var Log = class {
    /**
     *  The provider connected to the log used to fetch additional details
     *  if necessary.
     */
    provider;
    /**
     *  The transaction hash of the transaction this log occurred in. Use the
     *  [[Log-getTransaction]] to get the [[TransactionResponse]].
     */
    transactionHash;
    /**
     *  The block hash of the block this log occurred in. Use the
     *  [[Log-getBlock]] to get the [[Block]].
     */
    blockHash;
    /**
     *  The block number of the block this log occurred in. It is preferred
     *  to use the [[Block-hash]] when fetching the related [[Block]],
     *  since in the case of an orphaned block, the block at that height may
     *  have changed.
     */
    blockNumber;
    /**
     *  If the **Log** represents a block that was removed due to an orphaned
     *  block, this will be true.
     *
     *  This can only happen within an orphan event listener.
     */
    removed;
    /**
     *  The address of the contract that emitted this log.
     */
    address;
    /**
     *  The data included in this log when it was emitted.
     */
    data;
    /**
     *  The indexed topics included in this log when it was emitted.
     *
     *  All topics are included in the bloom filters, so they can be
     *  efficiently filtered using the [[Provider-getLogs]] method.
     */
    topics;
    /**
     *  The index within the block this log occurred at. This is generally
     *  not useful to developers, but can be used with the various roots
     *  to proof inclusion within a block.
     */
    index;
    /**
     *  The index within the transaction of this log.
     */
    transactionIndex;
    /**
     *  @_ignore:
     */
    constructor(log, provider) {
      this.provider = provider;
      const topics = Object.freeze(log.topics.slice());
      defineProperties(this, {
        transactionHash: log.transactionHash,
        blockHash: log.blockHash,
        blockNumber: log.blockNumber,
        removed: log.removed,
        address: log.address,
        data: log.data,
        topics,
        index: log.index,
        transactionIndex: log.transactionIndex
      });
    }
    /**
     *  Returns a JSON-compatible object.
     */
    toJSON() {
      const { address, blockHash, blockNumber, data, index, removed, topics, transactionHash, transactionIndex } = this;
      return {
        _type: "log",
        address,
        blockHash,
        blockNumber,
        data,
        index,
        removed,
        topics,
        transactionHash,
        transactionIndex
      };
    }
    /**
     *  Returns the block that this log occurred in.
     */
    async getBlock() {
      const block = await this.provider.getBlock(this.blockHash);
      assert(!!block, "failed to find transaction", "UNKNOWN_ERROR", {});
      return block;
    }
    /**
     *  Returns the transaction that this log occurred in.
     */
    async getTransaction() {
      const tx = await this.provider.getTransaction(this.transactionHash);
      assert(!!tx, "failed to find transaction", "UNKNOWN_ERROR", {});
      return tx;
    }
    /**
     *  Returns the transaction receipt fot the transaction that this
     *  log occurred in.
     */
    async getTransactionReceipt() {
      const receipt = await this.provider.getTransactionReceipt(this.transactionHash);
      assert(!!receipt, "failed to find transaction receipt", "UNKNOWN_ERROR", {});
      return receipt;
    }
    /**
     *  @_ignore:
     */
    removedEvent() {
      return createRemovedLogFilter(this);
    }
  };
  var TransactionReceipt = class {
    /**
     *  The provider connected to the log used to fetch additional details
     *  if necessary.
     */
    provider;
    /**
     *  The address the transaction was send to.
     */
    to;
    /**
     *  The sender of the transaction.
     */
    from;
    /**
     *  The address of the contract if the transaction was directly
     *  responsible for deploying one.
     *
     *  This is non-null **only** if the ``to`` is empty and the ``data``
     *  was successfully executed as initcode.
     */
    contractAddress;
    /**
     *  The transaction hash.
     */
    hash;
    /**
     *  The index of this transaction within the block transactions.
     */
    index;
    /**
     *  The block hash of the [[Block]] this transaction was included in.
     */
    blockHash;
    /**
     *  The block number of the [[Block]] this transaction was included in.
     */
    blockNumber;
    /**
     *  The bloom filter bytes that represent all logs that occurred within
     *  this transaction. This is generally not useful for most developers,
     *  but can be used to validate the included logs.
     */
    logsBloom;
    /**
     *  The actual amount of gas used by this transaction.
     *
     *  When creating a transaction, the amount of gas that will be used can
     *  only be approximated, but the sender must pay the gas fee for the
     *  entire gas limit. After the transaction, the difference is refunded.
     */
    gasUsed;
    /**
     *  The amount of gas used by all transactions within the block for this
     *  and all transactions with a lower ``index``.
     *
     *  This is generally not useful for developers but can be used to
     *  validate certain aspects of execution.
     */
    cumulativeGasUsed;
    /**
     *  The actual gas price used during execution.
     *
     *  Due to the complexity of [[link-eip-1559]] this value can only
     *  be caluclated after the transaction has been mined, snce the base
     *  fee is protocol-enforced.
     */
    gasPrice;
    /**
     *  The [[link-eip-2718]] transaction type.
     */
    type;
    //readonly byzantium!: boolean;
    /**
     *  The status of this transaction, indicating success (i.e. ``1``) or
     *  a revert (i.e. ``0``).
     *
     *  This is available in post-byzantium blocks, but some backends may
     *  backfill this value.
     */
    status;
    /**
     *  The root hash of this transaction.
     *
     *  This is no present and was only included in pre-byzantium blocks, but
     *  could be used to validate certain parts of the receipt.
     */
    root;
    #logs;
    /**
     *  @_ignore:
     */
    constructor(tx, provider) {
      this.#logs = Object.freeze(tx.logs.map((log) => {
        return new Log(log, provider);
      }));
      let gasPrice = BN_09;
      if (tx.effectiveGasPrice != null) {
        gasPrice = tx.effectiveGasPrice;
      } else if (tx.gasPrice != null) {
        gasPrice = tx.gasPrice;
      }
      defineProperties(this, {
        provider,
        to: tx.to,
        from: tx.from,
        contractAddress: tx.contractAddress,
        hash: tx.hash,
        index: tx.index,
        blockHash: tx.blockHash,
        blockNumber: tx.blockNumber,
        logsBloom: tx.logsBloom,
        gasUsed: tx.gasUsed,
        cumulativeGasUsed: tx.cumulativeGasUsed,
        gasPrice,
        type: tx.type,
        //byzantium: tx.byzantium,
        status: tx.status,
        root: tx.root
      });
    }
    /**
     *  The logs for this transaction.
     */
    get logs() {
      return this.#logs;
    }
    /**
     *  Returns a JSON-compatible representation.
     */
    toJSON() {
      const {
        to,
        from,
        contractAddress,
        hash: hash2,
        index,
        blockHash,
        blockNumber,
        logsBloom,
        logs,
        //byzantium, 
        status,
        root
      } = this;
      return {
        _type: "TransactionReceipt",
        blockHash,
        blockNumber,
        //byzantium, 
        contractAddress,
        cumulativeGasUsed: toJson(this.cumulativeGasUsed),
        from,
        gasPrice: toJson(this.gasPrice),
        gasUsed: toJson(this.gasUsed),
        hash: hash2,
        index,
        logs,
        logsBloom,
        root,
        status,
        to
      };
    }
    /**
     *  @_ignore:
     */
    get length() {
      return this.logs.length;
    }
    [Symbol.iterator]() {
      let index = 0;
      return {
        next: () => {
          if (index < this.length) {
            return { value: this.logs[index++], done: false };
          }
          return { value: void 0, done: true };
        }
      };
    }
    /**
     *  The total fee for this transaction, in wei.
     */
    get fee() {
      return this.gasUsed * this.gasPrice;
    }
    /**
     *  Resolves to the block this transaction occurred in.
     */
    async getBlock() {
      const block = await this.provider.getBlock(this.blockHash);
      if (block == null) {
        throw new Error("TODO");
      }
      return block;
    }
    /**
     *  Resolves to the transaction this transaction occurred in.
     */
    async getTransaction() {
      const tx = await this.provider.getTransaction(this.hash);
      if (tx == null) {
        throw new Error("TODO");
      }
      return tx;
    }
    /**
     *  Resolves to the return value of the execution of this transaction.
     *
     *  Support for this feature is limited, as it requires an archive node
     *  with the ``debug_`` or ``trace_`` API enabled.
     */
    async getResult() {
      return await this.provider.getTransactionResult(this.hash);
    }
    /**
     *  Resolves to the number of confirmations this transaction has.
     */
    async confirmations() {
      return await this.provider.getBlockNumber() - this.blockNumber + 1;
    }
    /**
     *  @_ignore:
     */
    removedEvent() {
      return createRemovedTransactionFilter(this);
    }
    /**
     *  @_ignore:
     */
    reorderedEvent(other) {
      assert(!other || other.isMined(), "unmined 'other' transction cannot be orphaned", "UNSUPPORTED_OPERATION", { operation: "reorderedEvent(other)" });
      return createReorderedTransactionFilter(this, other);
    }
  };
  var TransactionResponse = class _TransactionResponse {
    /**
     *  The provider this is connected to, which will influence how its
     *  methods will resolve its async inspection methods.
     */
    provider;
    /**
     *  The block number of the block that this transaction was included in.
     *
     *  This is ``null`` for pending transactions.
     */
    blockNumber;
    /**
     *  The blockHash of the block that this transaction was included in.
     *
     *  This is ``null`` for pending transactions.
     */
    blockHash;
    /**
     *  The index within the block that this transaction resides at.
     */
    index;
    /**
     *  The transaction hash.
     */
    hash;
    /**
     *  The [[link-eip-2718]] transaction envelope type. This is
     *  ``0`` for legacy transactions types.
     */
    type;
    /**
     *  The receiver of this transaction.
     *
     *  If ``null``, then the transaction is an initcode transaction.
     *  This means the result of executing the [[data]] will be deployed
     *  as a new contract on chain (assuming it does not revert) and the
     *  address may be computed using [[getCreateAddress]].
     */
    to;
    /**
     *  The sender of this transaction. It is implicitly computed
     *  from the transaction pre-image hash (as the digest) and the
     *  [[signature]] using ecrecover.
     */
    from;
    /**
     *  The nonce, which is used to prevent replay attacks and offer
     *  a method to ensure transactions from a given sender are explicitly
     *  ordered.
     *
     *  When sending a transaction, this must be equal to the number of
     *  transactions ever sent by [[from]].
     */
    nonce;
    /**
     *  The maximum units of gas this transaction can consume. If execution
     *  exceeds this, the entries transaction is reverted and the sender
     *  is charged for the full amount, despite not state changes being made.
     */
    gasLimit;
    /**
     *  The gas price can have various values, depending on the network.
     *
     *  In modern networks, for transactions that are included this is
     *  the //effective gas price// (the fee per gas that was actually
     *  charged), while for transactions that have not been included yet
     *  is the [[maxFeePerGas]].
     *
     *  For legacy transactions, or transactions on legacy networks, this
     *  is the fee that will be charged per unit of gas the transaction
     *  consumes.
     */
    gasPrice;
    /**
     *  The maximum priority fee (per unit of gas) to allow a
     *  validator to charge the sender. This is inclusive of the
     *  [[maxFeeFeePerGas]].
     */
    maxPriorityFeePerGas;
    /**
     *  The maximum fee (per unit of gas) to allow this transaction
     *  to charge the sender.
     */
    maxFeePerGas;
    /**
     *  The data.
     */
    data;
    /**
     *  The value, in wei. Use [[formatEther]] to format this value
     *  as ether.
     */
    value;
    /**
     *  The chain ID.
     */
    chainId;
    /**
     *  The signature.
     */
    signature;
    /**
     *  The [[link-eip-2930]] access list for transaction types that
     *  support it, otherwise ``null``.
     */
    accessList;
    #startBlock;
    /**
     *  @_ignore:
     */
    constructor(tx, provider) {
      this.provider = provider;
      this.blockNumber = tx.blockNumber != null ? tx.blockNumber : null;
      this.blockHash = tx.blockHash != null ? tx.blockHash : null;
      this.hash = tx.hash;
      this.index = tx.index;
      this.type = tx.type;
      this.from = tx.from;
      this.to = tx.to || null;
      this.gasLimit = tx.gasLimit;
      this.nonce = tx.nonce;
      this.data = tx.data;
      this.value = tx.value;
      this.gasPrice = tx.gasPrice;
      this.maxPriorityFeePerGas = tx.maxPriorityFeePerGas != null ? tx.maxPriorityFeePerGas : null;
      this.maxFeePerGas = tx.maxFeePerGas != null ? tx.maxFeePerGas : null;
      this.chainId = tx.chainId;
      this.signature = tx.signature;
      this.accessList = tx.accessList != null ? tx.accessList : null;
      this.#startBlock = -1;
    }
    /**
     *  Returns a JSON-compatible representation of this transaction.
     */
    toJSON() {
      const { blockNumber, blockHash, index, hash: hash2, type, to, from, nonce, data, signature, accessList } = this;
      return {
        _type: "TransactionReceipt",
        accessList,
        blockNumber,
        blockHash,
        chainId: toJson(this.chainId),
        data,
        from,
        gasLimit: toJson(this.gasLimit),
        gasPrice: toJson(this.gasPrice),
        hash: hash2,
        maxFeePerGas: toJson(this.maxFeePerGas),
        maxPriorityFeePerGas: toJson(this.maxPriorityFeePerGas),
        nonce,
        signature,
        to,
        index,
        type,
        value: toJson(this.value)
      };
    }
    /**
     *  Resolves to the Block that this transaction was included in.
     *
     *  This will return null if the transaction has not been included yet.
     */
    async getBlock() {
      let blockNumber = this.blockNumber;
      if (blockNumber == null) {
        const tx = await this.getTransaction();
        if (tx) {
          blockNumber = tx.blockNumber;
        }
      }
      if (blockNumber == null) {
        return null;
      }
      const block = this.provider.getBlock(blockNumber);
      if (block == null) {
        throw new Error("TODO");
      }
      return block;
    }
    /**
     *  Resolves to this transaction being re-requested from the
     *  provider. This can be used if you have an unmined transaction
     *  and wish to get an up-to-date populated instance.
     */
    async getTransaction() {
      return this.provider.getTransaction(this.hash);
    }
    /**
     *  Resolve to the number of confirmations this transaction has.
     */
    async confirmations() {
      if (this.blockNumber == null) {
        const { tx, blockNumber: blockNumber2 } = await resolveProperties({
          tx: this.getTransaction(),
          blockNumber: this.provider.getBlockNumber()
        });
        if (tx == null || tx.blockNumber == null) {
          return 0;
        }
        return blockNumber2 - tx.blockNumber + 1;
      }
      const blockNumber = await this.provider.getBlockNumber();
      return blockNumber - this.blockNumber + 1;
    }
    /**
     *  Resolves once this transaction has been mined and has
     *  %%confirms%% blocks including it (default: ``1``) with an
     *  optional %%timeout%%.
     *
     *  This can resolve to ``null`` only if %%confirms%% is ``0``
     *  and the transaction has not been mined, otherwise this will
     *  wait until enough confirmations have completed.
     */
    async wait(_confirms, _timeout) {
      const confirms = _confirms == null ? 1 : _confirms;
      const timeout = _timeout == null ? 0 : _timeout;
      let startBlock = this.#startBlock;
      let nextScan = -1;
      let stopScanning = startBlock === -1 ? true : false;
      const checkReplacement = async () => {
        if (stopScanning) {
          return null;
        }
        const { blockNumber, nonce } = await resolveProperties({
          blockNumber: this.provider.getBlockNumber(),
          nonce: this.provider.getTransactionCount(this.from)
        });
        if (nonce < this.nonce) {
          startBlock = blockNumber;
          return;
        }
        if (stopScanning) {
          return null;
        }
        const mined = await this.getTransaction();
        if (mined && mined.blockNumber != null) {
          return;
        }
        if (nextScan === -1) {
          nextScan = startBlock - 3;
          if (nextScan < this.#startBlock) {
            nextScan = this.#startBlock;
          }
        }
        while (nextScan <= blockNumber) {
          if (stopScanning) {
            return null;
          }
          const block = await this.provider.getBlock(nextScan, true);
          if (block == null) {
            return;
          }
          for (const hash2 of block) {
            if (hash2 === this.hash) {
              return;
            }
          }
          for (let i4 = 0; i4 < block.length; i4++) {
            const tx = await block.getTransaction(i4);
            if (tx.from === this.from && tx.nonce === this.nonce) {
              if (stopScanning) {
                return null;
              }
              const receipt2 = await this.provider.getTransactionReceipt(tx.hash);
              if (receipt2 == null) {
                return;
              }
              if (blockNumber - receipt2.blockNumber + 1 < confirms) {
                return;
              }
              let reason = "replaced";
              if (tx.data === this.data && tx.to === this.to && tx.value === this.value) {
                reason = "repriced";
              } else if (tx.data === "0x" && tx.from === tx.to && tx.value === BN_09) {
                reason = "cancelled";
              }
              assert(false, "transaction was replaced", "TRANSACTION_REPLACED", {
                cancelled: reason === "replaced" || reason === "cancelled",
                reason,
                replacement: tx.replaceableTransaction(startBlock),
                hash: tx.hash,
                receipt: receipt2
              });
            }
          }
          nextScan++;
        }
        return;
      };
      const checkReceipt = (receipt2) => {
        if (receipt2 == null || receipt2.status !== 0) {
          return receipt2;
        }
        assert(false, "transaction execution reverted", "CALL_EXCEPTION", {
          action: "sendTransaction",
          data: null,
          reason: null,
          invocation: null,
          revert: null,
          transaction: {
            to: receipt2.to,
            from: receipt2.from,
            data: ""
            // @TODO: in v7, split out sendTransaction properties
          },
          receipt: receipt2
        });
      };
      const receipt = await this.provider.getTransactionReceipt(this.hash);
      if (confirms === 0) {
        return checkReceipt(receipt);
      }
      if (receipt) {
        if (await receipt.confirmations() >= confirms) {
          return checkReceipt(receipt);
        }
      } else {
        await checkReplacement();
        if (confirms === 0) {
          return null;
        }
      }
      const waiter = new Promise((resolve, reject) => {
        const cancellers = [];
        const cancel = () => {
          cancellers.forEach((c4) => c4());
        };
        cancellers.push(() => {
          stopScanning = true;
        });
        if (timeout > 0) {
          const timer = setTimeout(() => {
            cancel();
            reject(makeError("wait for transaction timeout", "TIMEOUT"));
          }, timeout);
          cancellers.push(() => {
            clearTimeout(timer);
          });
        }
        const txListener = async (receipt2) => {
          if (await receipt2.confirmations() >= confirms) {
            cancel();
            try {
              resolve(checkReceipt(receipt2));
            } catch (error) {
              reject(error);
            }
          }
        };
        cancellers.push(() => {
          this.provider.off(this.hash, txListener);
        });
        this.provider.on(this.hash, txListener);
        if (startBlock >= 0) {
          const replaceListener = async () => {
            try {
              await checkReplacement();
            } catch (error) {
              if (isError(error, "TRANSACTION_REPLACED")) {
                cancel();
                reject(error);
                return;
              }
            }
            if (!stopScanning) {
              this.provider.once("block", replaceListener);
            }
          };
          cancellers.push(() => {
            this.provider.off("block", replaceListener);
          });
          this.provider.once("block", replaceListener);
        }
      });
      return await waiter;
    }
    /**
     *  Returns ``true`` if this transaction has been included.
     *
     *  This is effective only as of the time the TransactionResponse
     *  was instantiated. To get up-to-date information, use
     *  [[getTransaction]].
     *
     *  This provides a Type Guard that this transaction will have
     *  non-null property values for properties that are null for
     *  unmined transactions.
     */
    isMined() {
      return this.blockHash != null;
    }
    /**
     *  Returns true if the transaction is a legacy (i.e. ``type == 0``)
     *  transaction.
     *
     *  This provides a Type Guard that this transaction will have
     *  the ``null``-ness for hardfork-specific properties set correctly.
     */
    isLegacy() {
      return this.type === 0;
    }
    /**
     *  Returns true if the transaction is a Berlin (i.e. ``type == 1``)
     *  transaction. See [[link-eip-2070]].
     *
     *  This provides a Type Guard that this transaction will have
     *  the ``null``-ness for hardfork-specific properties set correctly.
     */
    isBerlin() {
      return this.type === 1;
    }
    /**
     *  Returns true if the transaction is a London (i.e. ``type == 2``)
     *  transaction. See [[link-eip-1559]].
     *
     *  This provides a Type Guard that this transaction will have
     *  the ``null``-ness for hardfork-specific properties set correctly.
     */
    isLondon() {
      return this.type === 2;
    }
    /**
     *  Returns a filter which can be used to listen for orphan events
     *  that evict this transaction.
     */
    removedEvent() {
      assert(this.isMined(), "unmined transaction canot be orphaned", "UNSUPPORTED_OPERATION", { operation: "removeEvent()" });
      return createRemovedTransactionFilter(this);
    }
    /**
     *  Returns a filter which can be used to listen for orphan events
     *  that re-order this event against %%other%%.
     */
    reorderedEvent(other) {
      assert(this.isMined(), "unmined transaction canot be orphaned", "UNSUPPORTED_OPERATION", { operation: "removeEvent()" });
      assert(!other || other.isMined(), "unmined 'other' transaction canot be orphaned", "UNSUPPORTED_OPERATION", { operation: "removeEvent()" });
      return createReorderedTransactionFilter(this, other);
    }
    /**
     *  Returns a new TransactionResponse instance which has the ability to
     *  detect (and throw an error) if the transaction is replaced, which
     *  will begin scanning at %%startBlock%%.
     *
     *  This should generally not be used by developers and is intended
     *  primarily for internal use. Setting an incorrect %%startBlock%% can
     *  have devastating performance consequences if used incorrectly.
     */
    replaceableTransaction(startBlock) {
      assertArgument(Number.isInteger(startBlock) && startBlock >= 0, "invalid startBlock", "startBlock", startBlock);
      const tx = new _TransactionResponse(this, this.provider);
      tx.#startBlock = startBlock;
      return tx;
    }
  };
  function createOrphanedBlockFilter(block) {
    return { orphan: "drop-block", hash: block.hash, number: block.number };
  }
  function createReorderedTransactionFilter(tx, other) {
    return { orphan: "reorder-transaction", tx, other };
  }
  function createRemovedTransactionFilter(tx) {
    return { orphan: "drop-transaction", tx };
  }
  function createRemovedLogFilter(log) {
    return { orphan: "drop-log", log: {
      transactionHash: log.transactionHash,
      blockHash: log.blockHash,
      blockNumber: log.blockNumber,
      address: log.address,
      data: log.data,
      topics: Object.freeze(log.topics.slice()),
      index: log.index
    } };
  }

  // node_modules/ethers/lib.esm/contract/wrappers.js
  var EventLog = class extends Log {
    /**
     *  The Contract Interface.
     */
    interface;
    /**
     *  The matching event.
     */
    fragment;
    /**
     *  The parsed arguments passed to the event by ``emit``.
     */
    args;
    /**
     * @_ignore:
     */
    constructor(log, iface, fragment) {
      super(log, log.provider);
      const args = iface.decodeEventLog(fragment, log.data, log.topics);
      defineProperties(this, { args, fragment, interface: iface });
    }
    /**
     *  The name of the event.
     */
    get eventName() {
      return this.fragment.name;
    }
    /**
     *  The signature of the event.
     */
    get eventSignature() {
      return this.fragment.format();
    }
  };
  var UndecodedEventLog = class extends Log {
    /**
     *  The error encounted when trying to decode the log.
     */
    error;
    /**
     * @_ignore:
     */
    constructor(log, error) {
      super(log, log.provider);
      defineProperties(this, { error });
    }
  };
  var ContractTransactionReceipt = class extends TransactionReceipt {
    #iface;
    /**
     *  @_ignore:
     */
    constructor(iface, provider, tx) {
      super(tx, provider);
      this.#iface = iface;
    }
    /**
     *  The parsed logs for any [[Log]] which has a matching event in the
     *  Contract ABI.
     */
    get logs() {
      return super.logs.map((log) => {
        const fragment = log.topics.length ? this.#iface.getEvent(log.topics[0]) : null;
        if (fragment) {
          try {
            return new EventLog(log, this.#iface, fragment);
          } catch (error) {
            return new UndecodedEventLog(log, error);
          }
        }
        return log;
      });
    }
  };
  var ContractTransactionResponse = class extends TransactionResponse {
    #iface;
    /**
     *  @_ignore:
     */
    constructor(iface, provider, tx) {
      super(tx, provider);
      this.#iface = iface;
    }
    /**
     *  Resolves once this transaction has been mined and has
     *  %%confirms%% blocks including it (default: ``1``) with an
     *  optional %%timeout%%.
     *
     *  This can resolve to ``null`` only if %%confirms%% is ``0``
     *  and the transaction has not been mined, otherwise this will
     *  wait until enough confirmations have completed.
     */
    async wait(confirms) {
      const receipt = await super.wait(confirms);
      if (receipt == null) {
        return null;
      }
      return new ContractTransactionReceipt(this.#iface, this.provider, receipt);
    }
  };
  var ContractUnknownEventPayload = class extends EventPayload {
    /**
     *  The log with no matching events.
     */
    log;
    /**
     *  @_event:
     */
    constructor(contract, listener, filter, log) {
      super(contract, listener, filter);
      defineProperties(this, { log });
    }
    /**
     *  Resolves to the block the event occured in.
     */
    async getBlock() {
      return await this.log.getBlock();
    }
    /**
     *  Resolves to the transaction the event occured in.
     */
    async getTransaction() {
      return await this.log.getTransaction();
    }
    /**
     *  Resolves to the transaction receipt the event occured in.
     */
    async getTransactionReceipt() {
      return await this.log.getTransactionReceipt();
    }
  };
  var ContractEventPayload = class extends ContractUnknownEventPayload {
    /**
     *  @_ignore:
     */
    constructor(contract, listener, filter, fragment, _log) {
      super(contract, listener, filter, new EventLog(_log, contract.interface, fragment));
      const args = contract.interface.decodeEventLog(fragment, this.log.data, this.log.topics);
      defineProperties(this, { args, fragment });
    }
    /**
     *  The event name.
     */
    get eventName() {
      return this.fragment.name;
    }
    /**
     *  The event signature.
     */
    get eventSignature() {
      return this.fragment.format();
    }
  };

  // node_modules/ethers/lib.esm/contract/contract.js
  var BN_010 = BigInt(0);
  function canCall(value) {
    return value && typeof value.call === "function";
  }
  function canEstimate(value) {
    return value && typeof value.estimateGas === "function";
  }
  function canResolve(value) {
    return value && typeof value.resolveName === "function";
  }
  function canSend(value) {
    return value && typeof value.sendTransaction === "function";
  }
  var PreparedTopicFilter = class {
    #filter;
    fragment;
    constructor(contract, fragment, args) {
      defineProperties(this, { fragment });
      if (fragment.inputs.length < args.length) {
        throw new Error("too many arguments");
      }
      const runner = getRunner(contract.runner, "resolveName");
      const resolver = canResolve(runner) ? runner : null;
      this.#filter = async function() {
        const resolvedArgs = await Promise.all(fragment.inputs.map((param, index) => {
          const arg = args[index];
          if (arg == null) {
            return null;
          }
          return param.walkAsync(args[index], (type, value) => {
            if (type === "address") {
              if (Array.isArray(value)) {
                return Promise.all(value.map((v3) => resolveAddress(v3, resolver)));
              }
              return resolveAddress(value, resolver);
            }
            return value;
          });
        }));
        return contract.interface.encodeFilterTopics(fragment, resolvedArgs);
      }();
    }
    getTopicFilter() {
      return this.#filter;
    }
  };
  function getRunner(value, feature) {
    if (value == null) {
      return null;
    }
    if (typeof value[feature] === "function") {
      return value;
    }
    if (value.provider && typeof value.provider[feature] === "function") {
      return value.provider;
    }
    return null;
  }
  function getProvider(value) {
    if (value == null) {
      return null;
    }
    return value.provider || null;
  }
  async function copyOverrides(arg, allowed) {
    const _overrides = Typed.dereference(arg, "overrides");
    assertArgument(typeof _overrides === "object", "invalid overrides parameter", "overrides", arg);
    const overrides = copyRequest(_overrides);
    assertArgument(overrides.to == null || (allowed || []).indexOf("to") >= 0, "cannot override to", "overrides.to", overrides.to);
    assertArgument(overrides.data == null || (allowed || []).indexOf("data") >= 0, "cannot override data", "overrides.data", overrides.data);
    if (overrides.from) {
      overrides.from = await resolveAddress(overrides.from);
    }
    return overrides;
  }
  async function resolveArgs(_runner, inputs, args) {
    const runner = getRunner(_runner, "resolveName");
    const resolver = canResolve(runner) ? runner : null;
    return await Promise.all(inputs.map((param, index) => {
      return param.walkAsync(args[index], (type, value) => {
        value = Typed.dereference(value, type);
        if (type === "address") {
          return resolveAddress(value, resolver);
        }
        return value;
      });
    }));
  }
  function buildWrappedFallback(contract) {
    const populateTransaction = async function(overrides) {
      const tx = await copyOverrides(overrides, ["data"]);
      tx.to = await contract.getAddress();
      const iface = contract.interface;
      const noValue = getBigInt(tx.value || BN_010, "overrides.value") === BN_010;
      const noData = (tx.data || "0x") === "0x";
      if (iface.fallback && !iface.fallback.payable && iface.receive && !noData && !noValue) {
        assertArgument(false, "cannot send data to receive or send value to non-payable fallback", "overrides", overrides);
      }
      assertArgument(iface.fallback || noData, "cannot send data to receive-only contract", "overrides.data", tx.data);
      const payable = iface.receive || iface.fallback && iface.fallback.payable;
      assertArgument(payable || noValue, "cannot send value to non-payable fallback", "overrides.value", tx.value);
      assertArgument(iface.fallback || noData, "cannot send data to receive-only contract", "overrides.data", tx.data);
      return tx;
    };
    const staticCall = async function(overrides) {
      const runner = getRunner(contract.runner, "call");
      assert(canCall(runner), "contract runner does not support calling", "UNSUPPORTED_OPERATION", { operation: "call" });
      const tx = await populateTransaction(overrides);
      try {
        return await runner.call(tx);
      } catch (error) {
        if (isCallException(error) && error.data) {
          throw contract.interface.makeError(error.data, tx);
        }
        throw error;
      }
    };
    const send = async function(overrides) {
      const runner = contract.runner;
      assert(canSend(runner), "contract runner does not support sending transactions", "UNSUPPORTED_OPERATION", { operation: "sendTransaction" });
      const tx = await runner.sendTransaction(await populateTransaction(overrides));
      const provider = getProvider(contract.runner);
      return new ContractTransactionResponse(contract.interface, provider, tx);
    };
    const estimateGas = async function(overrides) {
      const runner = getRunner(contract.runner, "estimateGas");
      assert(canEstimate(runner), "contract runner does not support gas estimation", "UNSUPPORTED_OPERATION", { operation: "estimateGas" });
      return await runner.estimateGas(await populateTransaction(overrides));
    };
    const method = async (overrides) => {
      return await send(overrides);
    };
    defineProperties(method, {
      _contract: contract,
      estimateGas,
      populateTransaction,
      send,
      staticCall
    });
    return method;
  }
  function buildWrappedMethod(contract, key) {
    const getFragment = function(...args) {
      const fragment = contract.interface.getFunction(key, args);
      assert(fragment, "no matching fragment", "UNSUPPORTED_OPERATION", {
        operation: "fragment",
        info: { key, args }
      });
      return fragment;
    };
    const populateTransaction = async function(...args) {
      const fragment = getFragment(...args);
      let overrides = {};
      if (fragment.inputs.length + 1 === args.length) {
        overrides = await copyOverrides(args.pop());
      }
      if (fragment.inputs.length !== args.length) {
        throw new Error("internal error: fragment inputs doesn't match arguments; should not happen");
      }
      const resolvedArgs = await resolveArgs(contract.runner, fragment.inputs, args);
      return Object.assign({}, overrides, await resolveProperties({
        to: contract.getAddress(),
        data: contract.interface.encodeFunctionData(fragment, resolvedArgs)
      }));
    };
    const staticCall = async function(...args) {
      const result = await staticCallResult(...args);
      if (result.length === 1) {
        return result[0];
      }
      return result;
    };
    const send = async function(...args) {
      const runner = contract.runner;
      assert(canSend(runner), "contract runner does not support sending transactions", "UNSUPPORTED_OPERATION", { operation: "sendTransaction" });
      const tx = await runner.sendTransaction(await populateTransaction(...args));
      const provider = getProvider(contract.runner);
      return new ContractTransactionResponse(contract.interface, provider, tx);
    };
    const estimateGas = async function(...args) {
      const runner = getRunner(contract.runner, "estimateGas");
      assert(canEstimate(runner), "contract runner does not support gas estimation", "UNSUPPORTED_OPERATION", { operation: "estimateGas" });
      return await runner.estimateGas(await populateTransaction(...args));
    };
    const staticCallResult = async function(...args) {
      const runner = getRunner(contract.runner, "call");
      assert(canCall(runner), "contract runner does not support calling", "UNSUPPORTED_OPERATION", { operation: "call" });
      const tx = await populateTransaction(...args);
      let result = "0x";
      try {
        result = await runner.call(tx);
      } catch (error) {
        if (isCallException(error) && error.data) {
          throw contract.interface.makeError(error.data, tx);
        }
        throw error;
      }
      const fragment = getFragment(...args);
      return contract.interface.decodeFunctionResult(fragment, result);
    };
    const method = async (...args) => {
      const fragment = getFragment(...args);
      if (fragment.constant) {
        return await staticCall(...args);
      }
      return await send(...args);
    };
    defineProperties(method, {
      name: contract.interface.getFunctionName(key),
      _contract: contract,
      _key: key,
      getFragment,
      estimateGas,
      populateTransaction,
      send,
      staticCall,
      staticCallResult
    });
    Object.defineProperty(method, "fragment", {
      configurable: false,
      enumerable: true,
      get: () => {
        const fragment = contract.interface.getFunction(key);
        assert(fragment, "no matching fragment", "UNSUPPORTED_OPERATION", {
          operation: "fragment",
          info: { key }
        });
        return fragment;
      }
    });
    return method;
  }
  function buildWrappedEvent(contract, key) {
    const getFragment = function(...args) {
      const fragment = contract.interface.getEvent(key, args);
      assert(fragment, "no matching fragment", "UNSUPPORTED_OPERATION", {
        operation: "fragment",
        info: { key, args }
      });
      return fragment;
    };
    const method = function(...args) {
      return new PreparedTopicFilter(contract, getFragment(...args), args);
    };
    defineProperties(method, {
      name: contract.interface.getEventName(key),
      _contract: contract,
      _key: key,
      getFragment
    });
    Object.defineProperty(method, "fragment", {
      configurable: false,
      enumerable: true,
      get: () => {
        const fragment = contract.interface.getEvent(key);
        assert(fragment, "no matching fragment", "UNSUPPORTED_OPERATION", {
          operation: "fragment",
          info: { key }
        });
        return fragment;
      }
    });
    return method;
  }
  var internal2 = Symbol.for("_ethersInternal_contract");
  var internalValues = /* @__PURE__ */ new WeakMap();
  function setInternal(contract, values) {
    internalValues.set(contract[internal2], values);
  }
  function getInternal(contract) {
    return internalValues.get(contract[internal2]);
  }
  function isDeferred(value) {
    return value && typeof value === "object" && "getTopicFilter" in value && typeof value.getTopicFilter === "function" && value.fragment;
  }
  async function getSubInfo(contract, event) {
    let topics;
    let fragment = null;
    if (Array.isArray(event)) {
      const topicHashify = function(name) {
        if (isHexString(name, 32)) {
          return name;
        }
        const fragment2 = contract.interface.getEvent(name);
        assertArgument(fragment2, "unknown fragment", "name", name);
        return fragment2.topicHash;
      };
      topics = event.map((e4) => {
        if (e4 == null) {
          return null;
        }
        if (Array.isArray(e4)) {
          return e4.map(topicHashify);
        }
        return topicHashify(e4);
      });
    } else if (event === "*") {
      topics = [null];
    } else if (typeof event === "string") {
      if (isHexString(event, 32)) {
        topics = [event];
      } else {
        fragment = contract.interface.getEvent(event);
        assertArgument(fragment, "unknown fragment", "event", event);
        topics = [fragment.topicHash];
      }
    } else if (isDeferred(event)) {
      topics = await event.getTopicFilter();
    } else if ("fragment" in event) {
      fragment = event.fragment;
      topics = [fragment.topicHash];
    } else {
      assertArgument(false, "unknown event name", "event", event);
    }
    topics = topics.map((t4) => {
      if (t4 == null) {
        return null;
      }
      if (Array.isArray(t4)) {
        const items = Array.from(new Set(t4.map((t5) => t5.toLowerCase())).values());
        if (items.length === 1) {
          return items[0];
        }
        items.sort();
        return items;
      }
      return t4.toLowerCase();
    });
    const tag = topics.map((t4) => {
      if (t4 == null) {
        return "null";
      }
      if (Array.isArray(t4)) {
        return t4.join("|");
      }
      return t4;
    }).join("&");
    return { fragment, tag, topics };
  }
  async function hasSub(contract, event) {
    const { subs } = getInternal(contract);
    return subs.get((await getSubInfo(contract, event)).tag) || null;
  }
  async function getSub(contract, operation, event) {
    const provider = getProvider(contract.runner);
    assert(provider, "contract runner does not support subscribing", "UNSUPPORTED_OPERATION", { operation });
    const { fragment, tag, topics } = await getSubInfo(contract, event);
    const { addr, subs } = getInternal(contract);
    let sub = subs.get(tag);
    if (!sub) {
      const address = addr ? addr : contract;
      const filter = { address, topics };
      const listener = (log) => {
        let foundFragment = fragment;
        if (foundFragment == null) {
          try {
            foundFragment = contract.interface.getEvent(log.topics[0]);
          } catch (error) {
          }
        }
        if (foundFragment) {
          const _foundFragment = foundFragment;
          const args = fragment ? contract.interface.decodeEventLog(fragment, log.data, log.topics) : [];
          emit(contract, event, args, (listener2) => {
            return new ContractEventPayload(contract, listener2, event, _foundFragment, log);
          });
        } else {
          emit(contract, event, [], (listener2) => {
            return new ContractUnknownEventPayload(contract, listener2, event, log);
          });
        }
      };
      let starting = [];
      const start = () => {
        if (starting.length) {
          return;
        }
        starting.push(provider.on(filter, listener));
      };
      const stop = async () => {
        if (starting.length == 0) {
          return;
        }
        let started = starting;
        starting = [];
        await Promise.all(started);
        provider.off(filter, listener);
      };
      sub = { tag, listeners: [], start, stop };
      subs.set(tag, sub);
    }
    return sub;
  }
  var lastEmit = Promise.resolve();
  async function _emit(contract, event, args, payloadFunc) {
    await lastEmit;
    const sub = await hasSub(contract, event);
    if (!sub) {
      return false;
    }
    const count = sub.listeners.length;
    sub.listeners = sub.listeners.filter(({ listener, once }) => {
      const passArgs = Array.from(args);
      if (payloadFunc) {
        passArgs.push(payloadFunc(once ? null : listener));
      }
      try {
        listener.call(contract, ...passArgs);
      } catch (error) {
      }
      return !once;
    });
    if (sub.listeners.length === 0) {
      sub.stop();
      getInternal(contract).subs.delete(sub.tag);
    }
    return count > 0;
  }
  async function emit(contract, event, args, payloadFunc) {
    try {
      await lastEmit;
    } catch (error) {
    }
    const resultPromise = _emit(contract, event, args, payloadFunc);
    lastEmit = resultPromise;
    return await resultPromise;
  }
  var passProperties2 = ["then"];
  var BaseContract = class _BaseContract {
    /**
     *  The target to connect to.
     *
     *  This can be an address, ENS name or any [[Addressable]], such as
     *  another contract. To get the resovled address, use the ``getAddress``
     *  method.
     */
    target;
    /**
     *  The contract Interface.
     */
    interface;
    /**
     *  The connected runner. This is generally a [[Provider]] or a
     *  [[Signer]], which dictates what operations are supported.
     *
     *  For example, a **Contract** connected to a [[Provider]] may
     *  only execute read-only operations.
     */
    runner;
    /**
     *  All the Events available on this contract.
     */
    filters;
    /**
     *  @_ignore:
     */
    [internal2];
    /**
     *  The fallback or receive function if any.
     */
    fallback;
    /**
     *  Creates a new contract connected to %%target%% with the %%abi%% and
     *  optionally connected to a %%runner%% to perform operations on behalf
     *  of.
     */
    constructor(target, abi, runner, _deployTx) {
      assertArgument(typeof target === "string" || isAddressable(target), "invalid value for Contract target", "target", target);
      if (runner == null) {
        runner = null;
      }
      const iface = Interface.from(abi);
      defineProperties(this, { target, runner, interface: iface });
      Object.defineProperty(this, internal2, { value: {} });
      let addrPromise;
      let addr = null;
      let deployTx = null;
      if (_deployTx) {
        const provider = getProvider(runner);
        deployTx = new ContractTransactionResponse(this.interface, provider, _deployTx);
      }
      let subs = /* @__PURE__ */ new Map();
      if (typeof target === "string") {
        if (isHexString(target)) {
          addr = target;
          addrPromise = Promise.resolve(target);
        } else {
          const resolver = getRunner(runner, "resolveName");
          if (!canResolve(resolver)) {
            throw makeError("contract runner does not support name resolution", "UNSUPPORTED_OPERATION", {
              operation: "resolveName"
            });
          }
          addrPromise = resolver.resolveName(target).then((addr2) => {
            if (addr2 == null) {
              throw makeError("an ENS name used for a contract target must be correctly configured", "UNCONFIGURED_NAME", {
                value: target
              });
            }
            getInternal(this).addr = addr2;
            return addr2;
          });
        }
      } else {
        addrPromise = target.getAddress().then((addr2) => {
          if (addr2 == null) {
            throw new Error("TODO");
          }
          getInternal(this).addr = addr2;
          return addr2;
        });
      }
      setInternal(this, { addrPromise, addr, deployTx, subs });
      const filters = new Proxy({}, {
        get: (target2, prop, receiver) => {
          if (typeof prop === "symbol" || passProperties2.indexOf(prop) >= 0) {
            return Reflect.get(target2, prop, receiver);
          }
          try {
            return this.getEvent(prop);
          } catch (error) {
            if (!isError(error, "INVALID_ARGUMENT") || error.argument !== "key") {
              throw error;
            }
          }
          return void 0;
        },
        has: (target2, prop) => {
          if (passProperties2.indexOf(prop) >= 0) {
            return Reflect.has(target2, prop);
          }
          return Reflect.has(target2, prop) || this.interface.hasEvent(String(prop));
        }
      });
      defineProperties(this, { filters });
      defineProperties(this, {
        fallback: iface.receive || iface.fallback ? buildWrappedFallback(this) : null
      });
      return new Proxy(this, {
        get: (target2, prop, receiver) => {
          if (typeof prop === "symbol" || prop in target2 || passProperties2.indexOf(prop) >= 0) {
            return Reflect.get(target2, prop, receiver);
          }
          try {
            return target2.getFunction(prop);
          } catch (error) {
            if (!isError(error, "INVALID_ARGUMENT") || error.argument !== "key") {
              throw error;
            }
          }
          return void 0;
        },
        has: (target2, prop) => {
          if (typeof prop === "symbol" || prop in target2 || passProperties2.indexOf(prop) >= 0) {
            return Reflect.has(target2, prop);
          }
          return target2.interface.hasFunction(prop);
        }
      });
    }
    /**
     *  Return a new Contract instance with the same target and ABI, but
     *  a different %%runner%%.
     */
    connect(runner) {
      return new _BaseContract(this.target, this.interface, runner);
    }
    /**
     *  Return a new Contract instance with the same ABI and runner, but
     *  a different %%target%%.
     */
    attach(target) {
      return new _BaseContract(target, this.interface, this.runner);
    }
    /**
     *  Return the resolved address of this Contract.
     */
    async getAddress() {
      return await getInternal(this).addrPromise;
    }
    /**
     *  Return the deployed bytecode or null if no bytecode is found.
     */
    async getDeployedCode() {
      const provider = getProvider(this.runner);
      assert(provider, "runner does not support .provider", "UNSUPPORTED_OPERATION", { operation: "getDeployedCode" });
      const code = await provider.getCode(await this.getAddress());
      if (code === "0x") {
        return null;
      }
      return code;
    }
    /**
     *  Resolve to this Contract once the bytecode has been deployed, or
     *  resolve immediately if already deployed.
     */
    async waitForDeployment() {
      const deployTx = this.deploymentTransaction();
      if (deployTx) {
        await deployTx.wait();
        return this;
      }
      const code = await this.getDeployedCode();
      if (code != null) {
        return this;
      }
      const provider = getProvider(this.runner);
      assert(provider != null, "contract runner does not support .provider", "UNSUPPORTED_OPERATION", { operation: "waitForDeployment" });
      return new Promise((resolve, reject) => {
        const checkCode = async () => {
          try {
            const code2 = await this.getDeployedCode();
            if (code2 != null) {
              return resolve(this);
            }
            provider.once("block", checkCode);
          } catch (error) {
            reject(error);
          }
        };
        checkCode();
      });
    }
    /**
     *  Return the transaction used to deploy this contract.
     *
     *  This is only available if this instance was returned from a
     *  [[ContractFactory]].
     */
    deploymentTransaction() {
      return getInternal(this).deployTx;
    }
    /**
     *  Return the function for a given name. This is useful when a contract
     *  method name conflicts with a JavaScript name such as ``prototype`` or
     *  when using a Contract programatically.
     */
    getFunction(key) {
      if (typeof key !== "string") {
        key = key.format();
      }
      const func = buildWrappedMethod(this, key);
      return func;
    }
    /**
     *  Return the event for a given name. This is useful when a contract
     *  event name conflicts with a JavaScript name such as ``prototype`` or
     *  when using a Contract programatically.
     */
    getEvent(key) {
      if (typeof key !== "string") {
        key = key.format();
      }
      return buildWrappedEvent(this, key);
    }
    /**
     *  @_ignore:
     */
    async queryTransaction(hash2) {
      throw new Error("@TODO");
    }
    /*
        // @TODO: this is a non-backwards compatible change, but will be added
        //        in v7 and in a potential SmartContract class in an upcoming
        //        v6 release
        async getTransactionReceipt(hash: string): Promise<null | ContractTransactionReceipt> {
            const provider = getProvider(this.runner);
            assert(provider, "contract runner does not have a provider",
                "UNSUPPORTED_OPERATION", { operation: "queryTransaction" });
    
            const receipt = await provider.getTransactionReceipt(hash);
            if (receipt == null) { return null; }
    
            return new ContractTransactionReceipt(this.interface, provider, receipt);
        }
        */
    /**
     *  Provide historic access to event data for %%event%% in the range
     *  %%fromBlock%% (default: ``0``) to %%toBlock%% (default: ``"latest"``)
     *  inclusive.
     */
    async queryFilter(event, fromBlock, toBlock) {
      if (fromBlock == null) {
        fromBlock = 0;
      }
      if (toBlock == null) {
        toBlock = "latest";
      }
      const { addr, addrPromise } = getInternal(this);
      const address = addr ? addr : await addrPromise;
      const { fragment, topics } = await getSubInfo(this, event);
      const filter = { address, topics, fromBlock, toBlock };
      const provider = getProvider(this.runner);
      assert(provider, "contract runner does not have a provider", "UNSUPPORTED_OPERATION", { operation: "queryFilter" });
      return (await provider.getLogs(filter)).map((log) => {
        let foundFragment = fragment;
        if (foundFragment == null) {
          try {
            foundFragment = this.interface.getEvent(log.topics[0]);
          } catch (error) {
          }
        }
        if (foundFragment) {
          try {
            return new EventLog(log, this.interface, foundFragment);
          } catch (error) {
            return new UndecodedEventLog(log, error);
          }
        }
        return new Log(log, provider);
      });
    }
    /**
     *  Add an event %%listener%% for the %%event%%.
     */
    async on(event, listener) {
      const sub = await getSub(this, "on", event);
      sub.listeners.push({ listener, once: false });
      sub.start();
      return this;
    }
    /**
     *  Add an event %%listener%% for the %%event%%, but remove the listener
     *  after it is fired once.
     */
    async once(event, listener) {
      const sub = await getSub(this, "once", event);
      sub.listeners.push({ listener, once: true });
      sub.start();
      return this;
    }
    /**
     *  Emit an %%event%% calling all listeners with %%args%%.
     *
     *  Resolves to ``true`` if any listeners were called.
     */
    async emit(event, ...args) {
      return await emit(this, event, args, null);
    }
    /**
     *  Resolves to the number of listeners of %%event%% or the total number
     *  of listeners if unspecified.
     */
    async listenerCount(event) {
      if (event) {
        const sub = await hasSub(this, event);
        if (!sub) {
          return 0;
        }
        return sub.listeners.length;
      }
      const { subs } = getInternal(this);
      let total = 0;
      for (const { listeners } of subs.values()) {
        total += listeners.length;
      }
      return total;
    }
    /**
     *  Resolves to the listeners subscribed to %%event%% or all listeners
     *  if unspecified.
     */
    async listeners(event) {
      if (event) {
        const sub = await hasSub(this, event);
        if (!sub) {
          return [];
        }
        return sub.listeners.map(({ listener }) => listener);
      }
      const { subs } = getInternal(this);
      let result = [];
      for (const { listeners } of subs.values()) {
        result = result.concat(listeners.map(({ listener }) => listener));
      }
      return result;
    }
    /**
     *  Remove the %%listener%% from the listeners for %%event%% or remove
     *  all listeners if unspecified.
     */
    async off(event, listener) {
      const sub = await hasSub(this, event);
      if (!sub) {
        return this;
      }
      if (listener) {
        const index = sub.listeners.map(({ listener: listener2 }) => listener2).indexOf(listener);
        if (index >= 0) {
          sub.listeners.splice(index, 1);
        }
      }
      if (listener == null || sub.listeners.length === 0) {
        sub.stop();
        getInternal(this).subs.delete(sub.tag);
      }
      return this;
    }
    /**
     *  Remove all the listeners for %%event%% or remove all listeners if
     *  unspecified.
     */
    async removeAllListeners(event) {
      if (event) {
        const sub = await hasSub(this, event);
        if (!sub) {
          return this;
        }
        sub.stop();
        getInternal(this).subs.delete(sub.tag);
      } else {
        const { subs } = getInternal(this);
        for (const { tag, stop } of subs.values()) {
          stop();
          subs.delete(tag);
        }
      }
      return this;
    }
    /**
     *  Alias for [on].
     */
    async addListener(event, listener) {
      return await this.on(event, listener);
    }
    /**
     *  Alias for [off].
     */
    async removeListener(event, listener) {
      return await this.off(event, listener);
    }
    /**
     *  Create a new Class for the %%abi%%.
     */
    static buildClass(abi) {
      class CustomContract extends _BaseContract {
        constructor(address, runner = null) {
          super(address, abi, runner);
        }
      }
      return CustomContract;
    }
    /**
     *  Create a new BaseContract with a specified Interface.
     */
    static from(target, abi, runner) {
      if (runner == null) {
        runner = null;
      }
      const contract = new this(target, abi, runner);
      return contract;
    }
  };
  function _ContractBase() {
    return BaseContract;
  }
  var Contract = class extends _ContractBase() {
  };

  // node_modules/ethers/lib.esm/providers/ens-resolver.js
  function getIpfsLink(link) {
    if (link.match(/^ipfs:\/\/ipfs\//i)) {
      link = link.substring(12);
    } else if (link.match(/^ipfs:\/\//i)) {
      link = link.substring(7);
    } else {
      assertArgument(false, "unsupported IPFS format", "link", link);
    }
    return `https://gateway.ipfs.io/ipfs/${link}`;
  }
  var MulticoinProviderPlugin = class {
    /**
     *  The name.
     */
    name;
    /**
     *  Creates a new **MulticoinProviderPluing** for %%name%%.
     */
    constructor(name) {
      defineProperties(this, { name });
    }
    connect(proivder) {
      return this;
    }
    /**
     *  Returns ``true`` if %%coinType%% is supported by this plugin.
     */
    supportsCoinType(coinType) {
      return false;
    }
    /**
     *  Resovles to the encoded %%address%% for %%coinType%%.
     */
    async encodeAddress(coinType, address) {
      throw new Error("unsupported coin");
    }
    /**
     *  Resovles to the decoded %%data%% for %%coinType%%.
     */
    async decodeAddress(coinType, data) {
      throw new Error("unsupported coin");
    }
  };
  var matcherIpfs = new RegExp("^(ipfs)://(.*)$", "i");
  var matchers = [
    new RegExp("^(https)://(.*)$", "i"),
    new RegExp("^(data):(.*)$", "i"),
    matcherIpfs,
    new RegExp("^eip155:[0-9]+/(erc[0-9]+):(.*)$", "i")
  ];
  var EnsResolver = class _EnsResolver {
    /**
     *  The connected provider.
     */
    provider;
    /**
     *  The address of the resolver.
     */
    address;
    /**
     *  The name this resolver was resolved against.
     */
    name;
    // For EIP-2544 names, the ancestor that provided the resolver
    #supports2544;
    #resolver;
    constructor(provider, address, name) {
      defineProperties(this, { provider, address, name });
      this.#supports2544 = null;
      this.#resolver = new Contract(address, [
        "function supportsInterface(bytes4) view returns (bool)",
        "function resolve(bytes, bytes) view returns (bytes)",
        "function addr(bytes32) view returns (address)",
        "function addr(bytes32, uint) view returns (bytes)",
        "function text(bytes32, string) view returns (string)",
        "function contenthash(bytes32) view returns (bytes)"
      ], provider);
    }
    /**
     *  Resolves to true if the resolver supports wildcard resolution.
     */
    async supportsWildcard() {
      if (this.#supports2544 == null) {
        this.#supports2544 = (async () => {
          try {
            return await this.#resolver.supportsInterface("0x9061b923");
          } catch (error) {
            if (isError(error, "CALL_EXCEPTION")) {
              return false;
            }
            this.#supports2544 = null;
            throw error;
          }
        })();
      }
      return await this.#supports2544;
    }
    async #fetch(funcName, params) {
      params = (params || []).slice();
      const iface = this.#resolver.interface;
      params.unshift(namehash(this.name));
      let fragment = null;
      if (await this.supportsWildcard()) {
        fragment = iface.getFunction(funcName);
        assert(fragment, "missing fragment", "UNKNOWN_ERROR", {
          info: { funcName }
        });
        params = [
          dnsEncode(this.name),
          iface.encodeFunctionData(fragment, params)
        ];
        funcName = "resolve(bytes,bytes)";
      }
      params.push({
        enableCcipRead: true
      });
      try {
        const result = await this.#resolver[funcName](...params);
        if (fragment) {
          return iface.decodeFunctionResult(fragment, result)[0];
        }
        return result;
      } catch (error) {
        if (!isError(error, "CALL_EXCEPTION")) {
          throw error;
        }
      }
      return null;
    }
    /**
     *  Resolves to the address for %%coinType%% or null if the
     *  provided %%coinType%% has not been configured.
     */
    async getAddress(coinType) {
      if (coinType == null) {
        coinType = 60;
      }
      if (coinType === 60) {
        try {
          const result = await this.#fetch("addr(bytes32)");
          if (result == null || result === ZeroAddress) {
            return null;
          }
          return result;
        } catch (error) {
          if (isError(error, "CALL_EXCEPTION")) {
            return null;
          }
          throw error;
        }
      }
      if (coinType >= 0 && coinType < 2147483648) {
        let ethCoinType = coinType + 2147483648;
        const data2 = await this.#fetch("addr(bytes32,uint)", [ethCoinType]);
        if (isHexString(data2, 20)) {
          return getAddress(data2);
        }
      }
      let coinPlugin = null;
      for (const plugin of this.provider.plugins) {
        if (!(plugin instanceof MulticoinProviderPlugin)) {
          continue;
        }
        if (plugin.supportsCoinType(coinType)) {
          coinPlugin = plugin;
          break;
        }
      }
      if (coinPlugin == null) {
        return null;
      }
      const data = await this.#fetch("addr(bytes32,uint)", [coinType]);
      if (data == null || data === "0x") {
        return null;
      }
      const address = await coinPlugin.decodeAddress(coinType, data);
      if (address != null) {
        return address;
      }
      assert(false, `invalid coin data`, "UNSUPPORTED_OPERATION", {
        operation: `getAddress(${coinType})`,
        info: { coinType, data }
      });
    }
    /**
     *  Resolves to the EIP-634 text record for %%key%%, or ``null``
     *  if unconfigured.
     */
    async getText(key) {
      const data = await this.#fetch("text(bytes32,string)", [key]);
      if (data == null || data === "0x") {
        return null;
      }
      return data;
    }
    /**
     *  Rsolves to the content-hash or ``null`` if unconfigured.
     */
    async getContentHash() {
      const data = await this.#fetch("contenthash(bytes32)");
      if (data == null || data === "0x") {
        return null;
      }
      const ipfs = data.match(/^0x(e3010170|e5010172)(([0-9a-f][0-9a-f])([0-9a-f][0-9a-f])([0-9a-f]*))$/);
      if (ipfs) {
        const scheme = ipfs[1] === "e3010170" ? "ipfs" : "ipns";
        const length = parseInt(ipfs[4], 16);
        if (ipfs[5].length === length * 2) {
          return `${scheme}://${encodeBase58("0x" + ipfs[2])}`;
        }
      }
      const swarm = data.match(/^0xe40101fa011b20([0-9a-f]*)$/);
      if (swarm && swarm[1].length === 64) {
        return `bzz://${swarm[1]}`;
      }
      assert(false, `invalid or unsupported content hash data`, "UNSUPPORTED_OPERATION", {
        operation: "getContentHash()",
        info: { data }
      });
    }
    /**
     *  Resolves to the avatar url or ``null`` if the avatar is either
     *  unconfigured or incorrectly configured (e.g. references an NFT
     *  not owned by the address).
     *
     *  If diagnosing issues with configurations, the [[_getAvatar]]
     *  method may be useful.
     */
    async getAvatar() {
      const avatar = await this._getAvatar();
      return avatar.url;
    }
    /**
     *  When resolving an avatar, there are many steps involved, such
     *  fetching metadata and possibly validating ownership of an
     *  NFT.
     *
     *  This method can be used to examine each step and the value it
     *  was working from.
     */
    async _getAvatar() {
      const linkage = [{ type: "name", value: this.name }];
      try {
        const avatar = await this.getText("avatar");
        if (avatar == null) {
          linkage.push({ type: "!avatar", value: "" });
          return { url: null, linkage };
        }
        linkage.push({ type: "avatar", value: avatar });
        for (let i4 = 0; i4 < matchers.length; i4++) {
          const match = avatar.match(matchers[i4]);
          if (match == null) {
            continue;
          }
          const scheme = match[1].toLowerCase();
          switch (scheme) {
            case "https":
            case "data":
              linkage.push({ type: "url", value: avatar });
              return { linkage, url: avatar };
            case "ipfs": {
              const url = getIpfsLink(avatar);
              linkage.push({ type: "ipfs", value: avatar });
              linkage.push({ type: "url", value: url });
              return { linkage, url };
            }
            case "erc721":
            case "erc1155": {
              const selector = scheme === "erc721" ? "tokenURI(uint256)" : "uri(uint256)";
              linkage.push({ type: scheme, value: avatar });
              const owner = await this.getAddress();
              if (owner == null) {
                linkage.push({ type: "!owner", value: "" });
                return { url: null, linkage };
              }
              const comps = (match[2] || "").split("/");
              if (comps.length !== 2) {
                linkage.push({ type: `!${scheme}caip`, value: match[2] || "" });
                return { url: null, linkage };
              }
              const tokenId = comps[1];
              const contract = new Contract(comps[0], [
                // ERC-721
                "function tokenURI(uint) view returns (string)",
                "function ownerOf(uint) view returns (address)",
                // ERC-1155
                "function uri(uint) view returns (string)",
                "function balanceOf(address, uint256) view returns (uint)"
              ], this.provider);
              if (scheme === "erc721") {
                const tokenOwner = await contract.ownerOf(tokenId);
                if (owner !== tokenOwner) {
                  linkage.push({ type: "!owner", value: tokenOwner });
                  return { url: null, linkage };
                }
                linkage.push({ type: "owner", value: tokenOwner });
              } else if (scheme === "erc1155") {
                const balance = await contract.balanceOf(owner, tokenId);
                if (!balance) {
                  linkage.push({ type: "!balance", value: "0" });
                  return { url: null, linkage };
                }
                linkage.push({ type: "balance", value: balance.toString() });
              }
              let metadataUrl = await contract[selector](tokenId);
              if (metadataUrl == null || metadataUrl === "0x") {
                linkage.push({ type: "!metadata-url", value: "" });
                return { url: null, linkage };
              }
              linkage.push({ type: "metadata-url-base", value: metadataUrl });
              if (scheme === "erc1155") {
                metadataUrl = metadataUrl.replace("{id}", toBeHex(tokenId, 32).substring(2));
                linkage.push({ type: "metadata-url-expanded", value: metadataUrl });
              }
              if (metadataUrl.match(/^ipfs:/i)) {
                metadataUrl = getIpfsLink(metadataUrl);
              }
              linkage.push({ type: "metadata-url", value: metadataUrl });
              let metadata = {};
              const response = await new FetchRequest(metadataUrl).send();
              response.assertOk();
              try {
                metadata = response.bodyJson;
              } catch (error) {
                try {
                  linkage.push({ type: "!metadata", value: response.bodyText });
                } catch (error2) {
                  const bytes2 = response.body;
                  if (bytes2) {
                    linkage.push({ type: "!metadata", value: hexlify(bytes2) });
                  }
                  return { url: null, linkage };
                }
                return { url: null, linkage };
              }
              if (!metadata) {
                linkage.push({ type: "!metadata", value: "" });
                return { url: null, linkage };
              }
              linkage.push({ type: "metadata", value: JSON.stringify(metadata) });
              let imageUrl = metadata.image;
              if (typeof imageUrl !== "string") {
                linkage.push({ type: "!imageUrl", value: "" });
                return { url: null, linkage };
              }
              if (imageUrl.match(/^(https:\/\/|data:)/i)) {
              } else {
                const ipfs = imageUrl.match(matcherIpfs);
                if (ipfs == null) {
                  linkage.push({ type: "!imageUrl-ipfs", value: imageUrl });
                  return { url: null, linkage };
                }
                linkage.push({ type: "imageUrl-ipfs", value: imageUrl });
                imageUrl = getIpfsLink(imageUrl);
              }
              linkage.push({ type: "url", value: imageUrl });
              return { linkage, url: imageUrl };
            }
          }
        }
      } catch (error) {
      }
      return { linkage, url: null };
    }
    static async getEnsAddress(provider) {
      const network = await provider.getNetwork();
      const ensPlugin = network.getPlugin("org.ethers.plugins.network.Ens");
      assert(ensPlugin, "network does not support ENS", "UNSUPPORTED_OPERATION", {
        operation: "getEnsAddress",
        info: { network }
      });
      return ensPlugin.address;
    }
    static async #getResolver(provider, name) {
      const ensAddr = await _EnsResolver.getEnsAddress(provider);
      try {
        const contract = new Contract(ensAddr, [
          "function resolver(bytes32) view returns (address)"
        ], provider);
        const addr = await contract.resolver(namehash(name), {
          enableCcipRead: true
        });
        if (addr === ZeroAddress) {
          return null;
        }
        return addr;
      } catch (error) {
        throw error;
      }
      return null;
    }
    /**
     *  Resolve to the ENS resolver for %%name%% using %%provider%% or
     *  ``null`` if unconfigured.
     */
    static async fromName(provider, name) {
      let currentName = name;
      while (true) {
        if (currentName === "" || currentName === ".") {
          return null;
        }
        if (name !== "eth" && currentName === "eth") {
          return null;
        }
        const addr = await _EnsResolver.#getResolver(provider, currentName);
        if (addr != null) {
          const resolver = new _EnsResolver(provider, addr, name);
          if (currentName !== name && !await resolver.supportsWildcard()) {
            return null;
          }
          return resolver;
        }
        currentName = currentName.split(".").slice(1).join(".");
      }
    }
  };

  // node_modules/ethers/lib.esm/providers/format.js
  var BN_011 = BigInt(0);
  function allowNull(format, nullValue) {
    return function(value) {
      if (value == null) {
        return nullValue;
      }
      return format(value);
    };
  }
  function arrayOf(format) {
    return (array) => {
      if (!Array.isArray(array)) {
        throw new Error("not an array");
      }
      return array.map((i4) => format(i4));
    };
  }
  function object(format, altNames) {
    return (value) => {
      const result = {};
      for (const key in format) {
        let srcKey = key;
        if (altNames && key in altNames && !(srcKey in value)) {
          for (const altKey of altNames[key]) {
            if (altKey in value) {
              srcKey = altKey;
              break;
            }
          }
        }
        try {
          const nv = format[key](value[srcKey]);
          if (nv !== void 0) {
            result[key] = nv;
          }
        } catch (error) {
          const message = error instanceof Error ? error.message : "not-an-error";
          assert(false, `invalid value for value.${key} (${message})`, "BAD_DATA", { value });
        }
      }
      return result;
    };
  }
  function formatBoolean(value) {
    switch (value) {
      case true:
      case "true":
        return true;
      case false:
      case "false":
        return false;
    }
    assertArgument(false, `invalid boolean; ${JSON.stringify(value)}`, "value", value);
  }
  function formatData(value) {
    assertArgument(isHexString(value, true), "invalid data", "value", value);
    return value;
  }
  function formatHash(value) {
    assertArgument(isHexString(value, 32), "invalid hash", "value", value);
    return value;
  }
  var _formatLog = object({
    address: getAddress,
    blockHash: formatHash,
    blockNumber: getNumber,
    data: formatData,
    index: getNumber,
    removed: allowNull(formatBoolean, false),
    topics: arrayOf(formatHash),
    transactionHash: formatHash,
    transactionIndex: getNumber
  }, {
    index: ["logIndex"]
  });
  function formatLog(value) {
    return _formatLog(value);
  }
  var _formatBlock = object({
    hash: allowNull(formatHash),
    parentHash: formatHash,
    number: getNumber,
    timestamp: getNumber,
    nonce: allowNull(formatData),
    difficulty: getBigInt,
    gasLimit: getBigInt,
    gasUsed: getBigInt,
    miner: allowNull(getAddress),
    extraData: formatData,
    baseFeePerGas: allowNull(getBigInt)
  });
  function formatBlock(value) {
    const result = _formatBlock(value);
    result.transactions = value.transactions.map((tx) => {
      if (typeof tx === "string") {
        return tx;
      }
      return formatTransactionResponse(tx);
    });
    return result;
  }
  var _formatReceiptLog = object({
    transactionIndex: getNumber,
    blockNumber: getNumber,
    transactionHash: formatHash,
    address: getAddress,
    topics: arrayOf(formatHash),
    data: formatData,
    index: getNumber,
    blockHash: formatHash
  }, {
    index: ["logIndex"]
  });
  function formatReceiptLog(value) {
    return _formatReceiptLog(value);
  }
  var _formatTransactionReceipt = object({
    to: allowNull(getAddress, null),
    from: allowNull(getAddress, null),
    contractAddress: allowNull(getAddress, null),
    // should be allowNull(hash), but broken-EIP-658 support is handled in receipt
    index: getNumber,
    root: allowNull(hexlify),
    gasUsed: getBigInt,
    logsBloom: allowNull(formatData),
    blockHash: formatHash,
    hash: formatHash,
    logs: arrayOf(formatReceiptLog),
    blockNumber: getNumber,
    //confirmations: allowNull(getNumber, null),
    cumulativeGasUsed: getBigInt,
    effectiveGasPrice: allowNull(getBigInt),
    status: allowNull(getNumber),
    type: allowNull(getNumber, 0)
  }, {
    effectiveGasPrice: ["gasPrice"],
    hash: ["transactionHash"],
    index: ["transactionIndex"]
  });
  function formatTransactionReceipt(value) {
    return _formatTransactionReceipt(value);
  }
  function formatTransactionResponse(value) {
    if (value.to && getBigInt(value.to) === BN_011) {
      value.to = "0x0000000000000000000000000000000000000000";
    }
    const result = object({
      hash: formatHash,
      type: (value2) => {
        if (value2 === "0x" || value2 == null) {
          return 0;
        }
        return getNumber(value2);
      },
      accessList: allowNull(accessListify, null),
      blockHash: allowNull(formatHash, null),
      blockNumber: allowNull(getNumber, null),
      transactionIndex: allowNull(getNumber, null),
      //confirmations: allowNull(getNumber, null),
      from: getAddress,
      // either (gasPrice) or (maxPriorityFeePerGas + maxFeePerGas) must be set
      gasPrice: allowNull(getBigInt),
      maxPriorityFeePerGas: allowNull(getBigInt),
      maxFeePerGas: allowNull(getBigInt),
      gasLimit: getBigInt,
      to: allowNull(getAddress, null),
      value: getBigInt,
      nonce: getNumber,
      data: formatData,
      creates: allowNull(getAddress, null),
      chainId: allowNull(getBigInt, null)
    }, {
      data: ["input"],
      gasLimit: ["gas"]
    })(value);
    if (result.to == null && result.creates == null) {
      result.creates = getCreateAddress(result);
    }
    if ((value.type === 1 || value.type === 2) && value.accessList == null) {
      result.accessList = [];
    }
    if (value.signature) {
      result.signature = Signature2.from(value.signature);
    } else {
      result.signature = Signature2.from(value);
    }
    if (result.chainId == null) {
      const chainId = result.signature.legacyChainId;
      if (chainId != null) {
        result.chainId = chainId;
      }
    }
    if (result.blockHash && getBigInt(result.blockHash) === BN_011) {
      result.blockHash = null;
    }
    return result;
  }

  // node_modules/ethers/lib.esm/providers/plugins-network.js
  var EnsAddress = "0x00000000000C2E074eC69A0dFb2997BA6C7d2e1e";
  var NetworkPlugin = class _NetworkPlugin {
    /**
     *  The name of the plugin.
     *
     *  It is recommended to use reverse-domain-notation, which permits
     *  unique names with a known authority as well as hierarchal entries.
     */
    name;
    /**
     *  Creates a new **NetworkPlugin**.
     */
    constructor(name) {
      defineProperties(this, { name });
    }
    /**
     *  Creates a copy of this plugin.
     */
    clone() {
      return new _NetworkPlugin(this.name);
    }
  };
  var GasCostPlugin = class _GasCostPlugin extends NetworkPlugin {
    /**
     *  The block number to treat these values as valid from.
     *
     *  This allows a hardfork to have updated values included as well as
     *  mulutiple hardforks to be supported.
     */
    effectiveBlock;
    /**
     *  The transactions base fee.
     */
    txBase;
    /**
     *  The fee for creating a new account.
     */
    txCreate;
    /**
     *  The fee per zero-byte in the data.
     */
    txDataZero;
    /**
     *  The fee per non-zero-byte in the data.
     */
    txDataNonzero;
    /**
     *  The fee per storage key in the [[link-eip-2930]] access list.
     */
    txAccessListStorageKey;
    /**
     *  The fee per address in the [[link-eip-2930]] access list.
     */
    txAccessListAddress;
    /**
     *  Creates a new GasCostPlugin from %%effectiveBlock%% until the
     *  latest block or another GasCostPlugin supercedes that block number,
     *  with the associated %%costs%%.
     */
    constructor(effectiveBlock, costs) {
      if (effectiveBlock == null) {
        effectiveBlock = 0;
      }
      super(`org.ethers.network.plugins.GasCost#${effectiveBlock || 0}`);
      const props = { effectiveBlock };
      function set(name, nullish) {
        let value = (costs || {})[name];
        if (value == null) {
          value = nullish;
        }
        assertArgument(typeof value === "number", `invalud value for ${name}`, "costs", costs);
        props[name] = value;
      }
      set("txBase", 21e3);
      set("txCreate", 32e3);
      set("txDataZero", 4);
      set("txDataNonzero", 16);
      set("txAccessListStorageKey", 1900);
      set("txAccessListAddress", 2400);
      defineProperties(this, props);
    }
    clone() {
      return new _GasCostPlugin(this.effectiveBlock, this);
    }
  };
  var EnsPlugin = class _EnsPlugin extends NetworkPlugin {
    /**
     *  The ENS Registrty Contract address.
     */
    address;
    /**
     *  The chain ID that the ENS contract lives on.
     */
    targetNetwork;
    /**
     *  Creates a new **EnsPlugin** connected to %%address%% on the
     *  %%targetNetwork%%. The default ENS address and mainnet is used
     *  if unspecified.
     */
    constructor(address, targetNetwork) {
      super("org.ethers.plugins.network.Ens");
      defineProperties(this, {
        address: address || EnsAddress,
        targetNetwork: targetNetwork == null ? 1 : targetNetwork
      });
    }
    clone() {
      return new _EnsPlugin(this.address, this.targetNetwork);
    }
  };
  var FetchUrlFeeDataNetworkPlugin = class extends NetworkPlugin {
    #url;
    #processFunc;
    /**
     *  The URL to initialize the FetchRequest with in %%processFunc%%.
     */
    get url() {
      return this.#url;
    }
    /**
     *  The callback to use when computing the FeeData.
     */
    get processFunc() {
      return this.#processFunc;
    }
    /**
     *  Creates a new **FetchUrlFeeDataNetworkPlugin** which will
     *  be used when computing the fee data for the network.
     */
    constructor(url, processFunc) {
      super("org.ethers.plugins.network.FetchUrlFeeDataPlugin");
      this.#url = url;
      this.#processFunc = processFunc;
    }
    // We are immutable, so we can serve as our own clone
    clone() {
      return this;
    }
  };

  // node_modules/ethers/lib.esm/providers/network.js
  var Networks = /* @__PURE__ */ new Map();
  var Network = class _Network {
    #name;
    #chainId;
    #plugins;
    /**
     *  Creates a new **Network** for %%name%% and %%chainId%%.
     */
    constructor(name, chainId) {
      this.#name = name;
      this.#chainId = getBigInt(chainId);
      this.#plugins = /* @__PURE__ */ new Map();
    }
    /**
     *  Returns a JSON-compatible representation of a Network.
     */
    toJSON() {
      return { name: this.name, chainId: String(this.chainId) };
    }
    /**
     *  The network common name.
     *
     *  This is the canonical name, as networks migh have multiple
     *  names.
     */
    get name() {
      return this.#name;
    }
    set name(value) {
      this.#name = value;
    }
    /**
     *  The network chain ID.
     */
    get chainId() {
      return this.#chainId;
    }
    set chainId(value) {
      this.#chainId = getBigInt(value, "chainId");
    }
    /**
     *  Returns true if %%other%% matches this network. Any chain ID
     *  must match, and if no chain ID is present, the name must match.
     *
     *  This method does not currently check for additional properties,
     *  such as ENS address or plug-in compatibility.
     */
    matches(other) {
      if (other == null) {
        return false;
      }
      if (typeof other === "string") {
        try {
          return this.chainId === getBigInt(other);
        } catch (error) {
        }
        return this.name === other;
      }
      if (typeof other === "number" || typeof other === "bigint") {
        try {
          return this.chainId === getBigInt(other);
        } catch (error) {
        }
        return false;
      }
      if (typeof other === "object") {
        if (other.chainId != null) {
          try {
            return this.chainId === getBigInt(other.chainId);
          } catch (error) {
          }
          return false;
        }
        if (other.name != null) {
          return this.name === other.name;
        }
        return false;
      }
      return false;
    }
    /**
     *  Returns the list of plugins currently attached to this Network.
     */
    get plugins() {
      return Array.from(this.#plugins.values());
    }
    /**
     *  Attach a new %%plugin%% to this Network. The network name
     *  must be unique, excluding any fragment.
     */
    attachPlugin(plugin) {
      if (this.#plugins.get(plugin.name)) {
        throw new Error(`cannot replace existing plugin: ${plugin.name} `);
      }
      this.#plugins.set(plugin.name, plugin.clone());
      return this;
    }
    /**
     *  Return the plugin, if any, matching %%name%% exactly. Plugins
     *  with fragments will not be returned unless %%name%% includes
     *  a fragment.
     */
    getPlugin(name) {
      return this.#plugins.get(name) || null;
    }
    /**
     *  Gets a list of all plugins that match %%name%%, with otr without
     *  a fragment.
     */
    getPlugins(basename) {
      return this.plugins.filter((p4) => p4.name.split("#")[0] === basename);
    }
    /**
     *  Create a copy of this Network.
     */
    clone() {
      const clone = new _Network(this.name, this.chainId);
      this.plugins.forEach((plugin) => {
        clone.attachPlugin(plugin.clone());
      });
      return clone;
    }
    /**
     *  Compute the intrinsic gas required for a transaction.
     *
     *  A GasCostPlugin can be attached to override the default
     *  values.
     */
    computeIntrinsicGas(tx) {
      const costs = this.getPlugin("org.ethers.plugins.network.GasCost") || new GasCostPlugin();
      let gas = costs.txBase;
      if (tx.to == null) {
        gas += costs.txCreate;
      }
      if (tx.data) {
        for (let i4 = 2; i4 < tx.data.length; i4 += 2) {
          if (tx.data.substring(i4, i4 + 2) === "00") {
            gas += costs.txDataZero;
          } else {
            gas += costs.txDataNonzero;
          }
        }
      }
      if (tx.accessList) {
        const accessList = accessListify(tx.accessList);
        for (const addr in accessList) {
          gas += costs.txAccessListAddress + costs.txAccessListStorageKey * accessList[addr].storageKeys.length;
        }
      }
      return gas;
    }
    /**
     *  Returns a new Network for the %%network%% name or chainId.
     */
    static from(network) {
      injectCommonNetworks();
      if (network == null) {
        return _Network.from("mainnet");
      }
      if (typeof network === "number") {
        network = BigInt(network);
      }
      if (typeof network === "string" || typeof network === "bigint") {
        const networkFunc = Networks.get(network);
        if (networkFunc) {
          return networkFunc();
        }
        if (typeof network === "bigint") {
          return new _Network("unknown", network);
        }
        assertArgument(false, "unknown network", "network", network);
      }
      if (typeof network.clone === "function") {
        const clone = network.clone();
        return clone;
      }
      if (typeof network === "object") {
        assertArgument(typeof network.name === "string" && typeof network.chainId === "number", "invalid network object name or chainId", "network", network);
        const custom = new _Network(network.name, network.chainId);
        if (network.ensAddress || network.ensNetwork != null) {
          custom.attachPlugin(new EnsPlugin(network.ensAddress, network.ensNetwork));
        }
        return custom;
      }
      assertArgument(false, "invalid network", "network", network);
    }
    /**
     *  Register %%nameOrChainId%% with a function which returns
     *  an instance of a Network representing that chain.
     */
    static register(nameOrChainId, networkFunc) {
      if (typeof nameOrChainId === "number") {
        nameOrChainId = BigInt(nameOrChainId);
      }
      const existing = Networks.get(nameOrChainId);
      if (existing) {
        assertArgument(false, `conflicting network for ${JSON.stringify(existing.name)}`, "nameOrChainId", nameOrChainId);
      }
      Networks.set(nameOrChainId, networkFunc);
    }
  };
  function parseUnits2(_value, decimals) {
    const value = String(_value);
    if (!value.match(/^[0-9.]+$/)) {
      throw new Error(`invalid gwei value: ${_value}`);
    }
    const comps = value.split(".");
    if (comps.length === 1) {
      comps.push("");
    }
    if (comps.length !== 2) {
      throw new Error(`invalid gwei value: ${_value}`);
    }
    while (comps[1].length < decimals) {
      comps[1] += "0";
    }
    if (comps[1].length > 9) {
      let frac = BigInt(comps[1].substring(0, 9));
      if (!comps[1].substring(9).match(/^0+$/)) {
        frac++;
      }
      comps[1] = frac.toString();
    }
    return BigInt(comps[0] + comps[1]);
  }
  function getGasStationPlugin(url) {
    return new FetchUrlFeeDataNetworkPlugin(url, async (fetchFeeData, provider, request) => {
      request.setHeader("User-Agent", "ethers");
      let response;
      try {
        const [_response, _feeData] = await Promise.all([
          request.send(),
          fetchFeeData()
        ]);
        response = _response;
        const payload = response.bodyJson.standard;
        const feeData = {
          gasPrice: _feeData.gasPrice,
          maxFeePerGas: parseUnits2(payload.maxFee, 9),
          maxPriorityFeePerGas: parseUnits2(payload.maxPriorityFee, 9)
        };
        return feeData;
      } catch (error) {
        assert(false, `error encountered with polygon gas station (${JSON.stringify(request.url)})`, "SERVER_ERROR", { request, response, error });
      }
    });
  }
  function getPriorityFeePlugin(maxPriorityFeePerGas) {
    return new FetchUrlFeeDataNetworkPlugin("data:", async (fetchFeeData, provider, request) => {
      const feeData = await fetchFeeData();
      if (feeData.maxFeePerGas == null || feeData.maxPriorityFeePerGas == null) {
        return feeData;
      }
      const baseFee = feeData.maxFeePerGas - feeData.maxPriorityFeePerGas;
      return {
        gasPrice: feeData.gasPrice,
        maxFeePerGas: baseFee + maxPriorityFeePerGas,
        maxPriorityFeePerGas
      };
    });
  }
  var injected = false;
  function injectCommonNetworks() {
    if (injected) {
      return;
    }
    injected = true;
    function registerEth(name, chainId, options) {
      const func = function() {
        const network = new Network(name, chainId);
        if (options.ensNetwork != null) {
          network.attachPlugin(new EnsPlugin(null, options.ensNetwork));
        }
        network.attachPlugin(new GasCostPlugin());
        (options.plugins || []).forEach((plugin) => {
          network.attachPlugin(plugin);
        });
        return network;
      };
      Network.register(name, func);
      Network.register(chainId, func);
      if (options.altNames) {
        options.altNames.forEach((name2) => {
          Network.register(name2, func);
        });
      }
    }
    registerEth("mainnet", 1, { ensNetwork: 1, altNames: ["homestead"] });
    registerEth("ropsten", 3, { ensNetwork: 3 });
    registerEth("rinkeby", 4, { ensNetwork: 4 });
    registerEth("goerli", 5, { ensNetwork: 5 });
    registerEth("kovan", 42, { ensNetwork: 42 });
    registerEth("sepolia", 11155111, {});
    registerEth("classic", 61, {});
    registerEth("classicKotti", 6, {});
    registerEth("arbitrum", 42161, {
      ensNetwork: 1
    });
    registerEth("arbitrum-goerli", 421613, {});
    registerEth("bnb", 56, { ensNetwork: 1 });
    registerEth("bnbt", 97, {});
    registerEth("linea", 59144, { ensNetwork: 1 });
    registerEth("linea-goerli", 59140, {});
    registerEth("matic", 137, {
      ensNetwork: 1,
      plugins: [
        getGasStationPlugin("https://gasstation.polygon.technology/v2")
      ]
    });
    registerEth("matic-mumbai", 80001, {
      altNames: ["maticMumbai", "maticmum"],
      plugins: [
        getGasStationPlugin("https://gasstation-testnet.polygon.technology/v2")
      ]
    });
    registerEth("optimism", 10, {
      ensNetwork: 1,
      plugins: [
        getPriorityFeePlugin(BigInt("1000000"))
      ]
    });
    registerEth("optimism-goerli", 420, {});
    registerEth("xdai", 100, { ensNetwork: 1 });
  }

  // node_modules/ethers/lib.esm/providers/subscriber-polling.js
  function copy(obj) {
    return JSON.parse(JSON.stringify(obj));
  }
  var PollingBlockSubscriber = class {
    #provider;
    #poller;
    #interval;
    // The most recent block we have scanned for events. The value -2
    // indicates we still need to fetch an initial block number
    #blockNumber;
    /**
     *  Create a new **PollingBlockSubscriber** attached to %%provider%%.
     */
    constructor(provider) {
      this.#provider = provider;
      this.#poller = null;
      this.#interval = 4e3;
      this.#blockNumber = -2;
    }
    /**
     *  The polling interval.
     */
    get pollingInterval() {
      return this.#interval;
    }
    set pollingInterval(value) {
      this.#interval = value;
    }
    async #poll() {
      try {
        const blockNumber = await this.#provider.getBlockNumber();
        if (this.#blockNumber === -2) {
          this.#blockNumber = blockNumber;
          return;
        }
        if (blockNumber !== this.#blockNumber) {
          for (let b5 = this.#blockNumber + 1; b5 <= blockNumber; b5++) {
            if (this.#poller == null) {
              return;
            }
            await this.#provider.emit("block", b5);
          }
          this.#blockNumber = blockNumber;
        }
      } catch (error) {
      }
      if (this.#poller == null) {
        return;
      }
      this.#poller = this.#provider._setTimeout(this.#poll.bind(this), this.#interval);
    }
    start() {
      if (this.#poller) {
        return;
      }
      this.#poller = this.#provider._setTimeout(this.#poll.bind(this), this.#interval);
      this.#poll();
    }
    stop() {
      if (!this.#poller) {
        return;
      }
      this.#provider._clearTimeout(this.#poller);
      this.#poller = null;
    }
    pause(dropWhilePaused) {
      this.stop();
      if (dropWhilePaused) {
        this.#blockNumber = -2;
      }
    }
    resume() {
      this.start();
    }
  };
  var OnBlockSubscriber = class {
    #provider;
    #poll;
    #running;
    /**
     *  Create a new **OnBlockSubscriber** attached to %%provider%%.
     */
    constructor(provider) {
      this.#provider = provider;
      this.#running = false;
      this.#poll = (blockNumber) => {
        this._poll(blockNumber, this.#provider);
      };
    }
    /**
     *  Called on every new block.
     */
    async _poll(blockNumber, provider) {
      throw new Error("sub-classes must override this");
    }
    start() {
      if (this.#running) {
        return;
      }
      this.#running = true;
      this.#poll(-2);
      this.#provider.on("block", this.#poll);
    }
    stop() {
      if (!this.#running) {
        return;
      }
      this.#running = false;
      this.#provider.off("block", this.#poll);
    }
    pause(dropWhilePaused) {
      this.stop();
    }
    resume() {
      this.start();
    }
  };
  var PollingOrphanSubscriber = class extends OnBlockSubscriber {
    #filter;
    constructor(provider, filter) {
      super(provider);
      this.#filter = copy(filter);
    }
    async _poll(blockNumber, provider) {
      throw new Error("@TODO");
      console.log(this.#filter);
    }
  };
  var PollingTransactionSubscriber = class extends OnBlockSubscriber {
    #hash;
    /**
     *  Create a new **PollingTransactionSubscriber** attached to
     *  %%provider%%, listening for %%hash%%.
     */
    constructor(provider, hash2) {
      super(provider);
      this.#hash = hash2;
    }
    async _poll(blockNumber, provider) {
      const tx = await provider.getTransactionReceipt(this.#hash);
      if (tx) {
        provider.emit(this.#hash, tx);
      }
    }
  };
  var PollingEventSubscriber = class {
    #provider;
    #filter;
    #poller;
    #running;
    // The most recent block we have scanned for events. The value -2
    // indicates we still need to fetch an initial block number
    #blockNumber;
    /**
     *  Create a new **PollingTransactionSubscriber** attached to
     *  %%provider%%, listening for %%filter%%.
     */
    constructor(provider, filter) {
      this.#provider = provider;
      this.#filter = copy(filter);
      this.#poller = this.#poll.bind(this);
      this.#running = false;
      this.#blockNumber = -2;
    }
    async #poll(blockNumber) {
      if (this.#blockNumber === -2) {
        return;
      }
      const filter = copy(this.#filter);
      filter.fromBlock = this.#blockNumber + 1;
      filter.toBlock = blockNumber;
      const logs = await this.#provider.getLogs(filter);
      if (logs.length === 0) {
        if (this.#blockNumber < blockNumber - 60) {
          this.#blockNumber = blockNumber - 60;
        }
        return;
      }
      for (const log of logs) {
        this.#provider.emit(this.#filter, log);
        this.#blockNumber = log.blockNumber;
      }
    }
    start() {
      if (this.#running) {
        return;
      }
      this.#running = true;
      if (this.#blockNumber === -2) {
        this.#provider.getBlockNumber().then((blockNumber) => {
          this.#blockNumber = blockNumber;
        });
      }
      this.#provider.on("block", this.#poller);
    }
    stop() {
      if (!this.#running) {
        return;
      }
      this.#running = false;
      this.#provider.off("block", this.#poller);
    }
    pause(dropWhilePaused) {
      this.stop();
      if (dropWhilePaused) {
        this.#blockNumber = -2;
      }
    }
    resume() {
      this.start();
    }
  };

  // node_modules/ethers/lib.esm/providers/abstract-provider.js
  var BN_23 = BigInt(2);
  var MAX_CCIP_REDIRECTS = 10;
  function isPromise(value) {
    return value && typeof value.then === "function";
  }
  function getTag(prefix, value) {
    return prefix + ":" + JSON.stringify(value, (k3, v3) => {
      if (v3 == null) {
        return "null";
      }
      if (typeof v3 === "bigint") {
        return `bigint:${v3.toString()}`;
      }
      if (typeof v3 === "string") {
        return v3.toLowerCase();
      }
      if (typeof v3 === "object" && !Array.isArray(v3)) {
        const keys = Object.keys(v3);
        keys.sort();
        return keys.reduce((accum, key) => {
          accum[key] = v3[key];
          return accum;
        }, {});
      }
      return v3;
    });
  }
  var UnmanagedSubscriber = class {
    /**
     *  The name fof the event.
     */
    name;
    /**
     *  Create a new UnmanagedSubscriber with %%name%%.
     */
    constructor(name) {
      defineProperties(this, { name });
    }
    start() {
    }
    stop() {
    }
    pause(dropWhilePaused) {
    }
    resume() {
    }
  };
  function copy2(value) {
    return JSON.parse(JSON.stringify(value));
  }
  function concisify(items) {
    items = Array.from(new Set(items).values());
    items.sort();
    return items;
  }
  async function getSubscription(_event, provider) {
    if (_event == null) {
      throw new Error("invalid event");
    }
    if (Array.isArray(_event)) {
      _event = { topics: _event };
    }
    if (typeof _event === "string") {
      switch (_event) {
        case "block":
        case "pending":
        case "debug":
        case "error":
        case "network": {
          return { type: _event, tag: _event };
        }
      }
    }
    if (isHexString(_event, 32)) {
      const hash2 = _event.toLowerCase();
      return { type: "transaction", tag: getTag("tx", { hash: hash2 }), hash: hash2 };
    }
    if (_event.orphan) {
      const event = _event;
      return { type: "orphan", tag: getTag("orphan", event), filter: copy2(event) };
    }
    if (_event.address || _event.topics) {
      const event = _event;
      const filter = {
        topics: (event.topics || []).map((t4) => {
          if (t4 == null) {
            return null;
          }
          if (Array.isArray(t4)) {
            return concisify(t4.map((t5) => t5.toLowerCase()));
          }
          return t4.toLowerCase();
        })
      };
      if (event.address) {
        const addresses = [];
        const promises = [];
        const addAddress = (addr) => {
          if (isHexString(addr)) {
            addresses.push(addr);
          } else {
            promises.push((async () => {
              addresses.push(await resolveAddress(addr, provider));
            })());
          }
        };
        if (Array.isArray(event.address)) {
          event.address.forEach(addAddress);
        } else {
          addAddress(event.address);
        }
        if (promises.length) {
          await Promise.all(promises);
        }
        filter.address = concisify(addresses.map((a4) => a4.toLowerCase()));
      }
      return { filter, tag: getTag("event", filter), type: "event" };
    }
    assertArgument(false, "unknown ProviderEvent", "event", _event);
  }
  function getTime2() {
    return (/* @__PURE__ */ new Date()).getTime();
  }
  var defaultOptions = {
    cacheTimeout: 250,
    pollingInterval: 4e3
  };
  var AbstractProvider = class {
    #subs;
    #plugins;
    // null=unpaused, true=paused+dropWhilePaused, false=paused
    #pausedState;
    #destroyed;
    #networkPromise;
    #anyNetwork;
    #performCache;
    // The most recent block number if running an event or -1 if no "block" event
    #lastBlockNumber;
    #nextTimer;
    #timers;
    #disableCcipRead;
    #options;
    /**
     *  Create a new **AbstractProvider** connected to %%network%%, or
     *  use the various network detection capabilities to discover the
     *  [[Network]] if necessary.
     */
    constructor(_network, options) {
      this.#options = Object.assign({}, defaultOptions, options || {});
      if (_network === "any") {
        this.#anyNetwork = true;
        this.#networkPromise = null;
      } else if (_network) {
        const network = Network.from(_network);
        this.#anyNetwork = false;
        this.#networkPromise = Promise.resolve(network);
        setTimeout(() => {
          this.emit("network", network, null);
        }, 0);
      } else {
        this.#anyNetwork = false;
        this.#networkPromise = null;
      }
      this.#lastBlockNumber = -1;
      this.#performCache = /* @__PURE__ */ new Map();
      this.#subs = /* @__PURE__ */ new Map();
      this.#plugins = /* @__PURE__ */ new Map();
      this.#pausedState = null;
      this.#destroyed = false;
      this.#nextTimer = 1;
      this.#timers = /* @__PURE__ */ new Map();
      this.#disableCcipRead = false;
    }
    get pollingInterval() {
      return this.#options.pollingInterval;
    }
    /**
     *  Returns ``this``, to allow an **AbstractProvider** to implement
     *  the [[ContractRunner]] interface.
     */
    get provider() {
      return this;
    }
    /**
     *  Returns all the registered plug-ins.
     */
    get plugins() {
      return Array.from(this.#plugins.values());
    }
    /**
     *  Attach a new plug-in.
     */
    attachPlugin(plugin) {
      if (this.#plugins.get(plugin.name)) {
        throw new Error(`cannot replace existing plugin: ${plugin.name} `);
      }
      this.#plugins.set(plugin.name, plugin.connect(this));
      return this;
    }
    /**
     *  Get a plugin by name.
     */
    getPlugin(name) {
      return this.#plugins.get(name) || null;
    }
    /**
     *  Prevent any CCIP-read operation, regardless of whether requested
     *  in a [[call]] using ``enableCcipRead``.
     */
    get disableCcipRead() {
      return this.#disableCcipRead;
    }
    set disableCcipRead(value) {
      this.#disableCcipRead = !!value;
    }
    // Shares multiple identical requests made during the same 250ms
    async #perform(req) {
      const timeout = this.#options.cacheTimeout;
      if (timeout < 0) {
        return await this._perform(req);
      }
      const tag = getTag(req.method, req);
      let perform = this.#performCache.get(tag);
      if (!perform) {
        perform = this._perform(req);
        this.#performCache.set(tag, perform);
        setTimeout(() => {
          if (this.#performCache.get(tag) === perform) {
            this.#performCache.delete(tag);
          }
        }, timeout);
      }
      return await perform;
    }
    /**
     *  Resolves to the data for executing the CCIP-read operations.
     */
    async ccipReadFetch(tx, calldata, urls) {
      if (this.disableCcipRead || urls.length === 0 || tx.to == null) {
        return null;
      }
      const sender = tx.to.toLowerCase();
      const data = calldata.toLowerCase();
      const errorMessages = [];
      for (let i4 = 0; i4 < urls.length; i4++) {
        const url = urls[i4];
        const href = url.replace("{sender}", sender).replace("{data}", data);
        const request = new FetchRequest(href);
        if (url.indexOf("{data}") === -1) {
          request.body = { data, sender };
        }
        this.emit("debug", { action: "sendCcipReadFetchRequest", request, index: i4, urls });
        let errorMessage = "unknown error";
        const resp = await request.send();
        try {
          const result = resp.bodyJson;
          if (result.data) {
            this.emit("debug", { action: "receiveCcipReadFetchResult", request, result });
            return result.data;
          }
          if (result.message) {
            errorMessage = result.message;
          }
          this.emit("debug", { action: "receiveCcipReadFetchError", request, result });
        } catch (error) {
        }
        assert(resp.statusCode < 400 || resp.statusCode >= 500, `response not found during CCIP fetch: ${errorMessage}`, "OFFCHAIN_FAULT", { reason: "404_MISSING_RESOURCE", transaction: tx, info: { url, errorMessage } });
        errorMessages.push(errorMessage);
      }
      assert(false, `error encountered during CCIP fetch: ${errorMessages.map((m4) => JSON.stringify(m4)).join(", ")}`, "OFFCHAIN_FAULT", {
        reason: "500_SERVER_ERROR",
        transaction: tx,
        info: { urls, errorMessages }
      });
    }
    /**
     *  Provides the opportunity for a sub-class to wrap a block before
     *  returning it, to add additional properties or an alternate
     *  sub-class of [[Block]].
     */
    _wrapBlock(value, network) {
      return new Block(formatBlock(value), this);
    }
    /**
     *  Provides the opportunity for a sub-class to wrap a log before
     *  returning it, to add additional properties or an alternate
     *  sub-class of [[Log]].
     */
    _wrapLog(value, network) {
      return new Log(formatLog(value), this);
    }
    /**
     *  Provides the opportunity for a sub-class to wrap a transaction
     *  receipt before returning it, to add additional properties or an
     *  alternate sub-class of [[TransactionReceipt]].
     */
    _wrapTransactionReceipt(value, network) {
      return new TransactionReceipt(formatTransactionReceipt(value), this);
    }
    /**
     *  Provides the opportunity for a sub-class to wrap a transaction
     *  response before returning it, to add additional properties or an
     *  alternate sub-class of [[TransactionResponse]].
     */
    _wrapTransactionResponse(tx, network) {
      return new TransactionResponse(formatTransactionResponse(tx), this);
    }
    /**
     *  Resolves to the Network, forcing a network detection using whatever
     *  technique the sub-class requires.
     *
     *  Sub-classes **must** override this.
     */
    _detectNetwork() {
      assert(false, "sub-classes must implement this", "UNSUPPORTED_OPERATION", {
        operation: "_detectNetwork"
      });
    }
    /**
     *  Sub-classes should use this to perform all built-in operations. All
     *  methods sanitizes and normalizes the values passed into this.
     *
     *  Sub-classes **must** override this.
     */
    async _perform(req) {
      assert(false, `unsupported method: ${req.method}`, "UNSUPPORTED_OPERATION", {
        operation: req.method,
        info: req
      });
    }
    // State
    async getBlockNumber() {
      const blockNumber = getNumber(await this.#perform({ method: "getBlockNumber" }), "%response");
      if (this.#lastBlockNumber >= 0) {
        this.#lastBlockNumber = blockNumber;
      }
      return blockNumber;
    }
    /**
     *  Returns or resolves to the address for %%address%%, resolving ENS
     *  names and [[Addressable]] objects and returning if already an
     *  address.
     */
    _getAddress(address) {
      return resolveAddress(address, this);
    }
    /**
     *  Returns or resolves to a valid block tag for %%blockTag%%, resolving
     *  negative values and returning if already a valid block tag.
     */
    _getBlockTag(blockTag) {
      if (blockTag == null) {
        return "latest";
      }
      switch (blockTag) {
        case "earliest":
          return "0x0";
        case "latest":
        case "pending":
        case "safe":
        case "finalized":
          return blockTag;
      }
      if (isHexString(blockTag)) {
        if (isHexString(blockTag, 32)) {
          return blockTag;
        }
        return toQuantity(blockTag);
      }
      if (typeof blockTag === "bigint") {
        blockTag = getNumber(blockTag, "blockTag");
      }
      if (typeof blockTag === "number") {
        if (blockTag >= 0) {
          return toQuantity(blockTag);
        }
        if (this.#lastBlockNumber >= 0) {
          return toQuantity(this.#lastBlockNumber + blockTag);
        }
        return this.getBlockNumber().then((b5) => toQuantity(b5 + blockTag));
      }
      assertArgument(false, "invalid blockTag", "blockTag", blockTag);
    }
    /**
     *  Returns or resolves to a filter for %%filter%%, resolving any ENS
     *  names or [[Addressable]] object and returning if already a valid
     *  filter.
     */
    _getFilter(filter) {
      const topics = (filter.topics || []).map((t4) => {
        if (t4 == null) {
          return null;
        }
        if (Array.isArray(t4)) {
          return concisify(t4.map((t5) => t5.toLowerCase()));
        }
        return t4.toLowerCase();
      });
      const blockHash = "blockHash" in filter ? filter.blockHash : void 0;
      const resolve = (_address, fromBlock2, toBlock2) => {
        let address2 = void 0;
        switch (_address.length) {
          case 0:
            break;
          case 1:
            address2 = _address[0];
            break;
          default:
            _address.sort();
            address2 = _address;
        }
        if (blockHash) {
          if (fromBlock2 != null || toBlock2 != null) {
            throw new Error("invalid filter");
          }
        }
        const filter2 = {};
        if (address2) {
          filter2.address = address2;
        }
        if (topics.length) {
          filter2.topics = topics;
        }
        if (fromBlock2) {
          filter2.fromBlock = fromBlock2;
        }
        if (toBlock2) {
          filter2.toBlock = toBlock2;
        }
        if (blockHash) {
          filter2.blockHash = blockHash;
        }
        return filter2;
      };
      let address = [];
      if (filter.address) {
        if (Array.isArray(filter.address)) {
          for (const addr of filter.address) {
            address.push(this._getAddress(addr));
          }
        } else {
          address.push(this._getAddress(filter.address));
        }
      }
      let fromBlock = void 0;
      if ("fromBlock" in filter) {
        fromBlock = this._getBlockTag(filter.fromBlock);
      }
      let toBlock = void 0;
      if ("toBlock" in filter) {
        toBlock = this._getBlockTag(filter.toBlock);
      }
      if (address.filter((a4) => typeof a4 !== "string").length || fromBlock != null && typeof fromBlock !== "string" || toBlock != null && typeof toBlock !== "string") {
        return Promise.all([Promise.all(address), fromBlock, toBlock]).then((result) => {
          return resolve(result[0], result[1], result[2]);
        });
      }
      return resolve(address, fromBlock, toBlock);
    }
    /**
     *  Returns or resovles to a transaction for %%request%%, resolving
     *  any ENS names or [[Addressable]] and returning if already a valid
     *  transaction.
     */
    _getTransactionRequest(_request) {
      const request = copyRequest(_request);
      const promises = [];
      ["to", "from"].forEach((key) => {
        if (request[key] == null) {
          return;
        }
        const addr = resolveAddress(request[key]);
        if (isPromise(addr)) {
          promises.push(async function() {
            request[key] = await addr;
          }());
        } else {
          request[key] = addr;
        }
      });
      if (request.blockTag != null) {
        const blockTag = this._getBlockTag(request.blockTag);
        if (isPromise(blockTag)) {
          promises.push(async function() {
            request.blockTag = await blockTag;
          }());
        } else {
          request.blockTag = blockTag;
        }
      }
      if (promises.length) {
        return async function() {
          await Promise.all(promises);
          return request;
        }();
      }
      return request;
    }
    async getNetwork() {
      if (this.#networkPromise == null) {
        const detectNetwork = this._detectNetwork().then((network) => {
          this.emit("network", network, null);
          return network;
        }, (error) => {
          if (this.#networkPromise === detectNetwork) {
            this.#networkPromise = null;
          }
          throw error;
        });
        this.#networkPromise = detectNetwork;
        return (await detectNetwork).clone();
      }
      const networkPromise = this.#networkPromise;
      const [expected, actual] = await Promise.all([
        networkPromise,
        this._detectNetwork()
        // The actual connected network
      ]);
      if (expected.chainId !== actual.chainId) {
        if (this.#anyNetwork) {
          this.emit("network", actual, expected);
          if (this.#networkPromise === networkPromise) {
            this.#networkPromise = Promise.resolve(actual);
          }
        } else {
          assert(false, `network changed: ${expected.chainId} => ${actual.chainId} `, "NETWORK_ERROR", {
            event: "changed"
          });
        }
      }
      return expected.clone();
    }
    async getFeeData() {
      const network = await this.getNetwork();
      const getFeeDataFunc = async () => {
        const { _block, gasPrice } = await resolveProperties({
          _block: this.#getBlock("latest", false),
          gasPrice: (async () => {
            try {
              const gasPrice2 = await this.#perform({ method: "getGasPrice" });
              return getBigInt(gasPrice2, "%response");
            } catch (error) {
            }
            return null;
          })()
        });
        let maxFeePerGas = null;
        let maxPriorityFeePerGas = null;
        const block = this._wrapBlock(_block, network);
        if (block && block.baseFeePerGas) {
          maxPriorityFeePerGas = BigInt("1000000000");
          maxFeePerGas = block.baseFeePerGas * BN_23 + maxPriorityFeePerGas;
        }
        return new FeeData(gasPrice, maxFeePerGas, maxPriorityFeePerGas);
      };
      const plugin = network.getPlugin("org.ethers.plugins.network.FetchUrlFeeDataPlugin");
      if (plugin) {
        const req = new FetchRequest(plugin.url);
        const feeData = await plugin.processFunc(getFeeDataFunc, this, req);
        return new FeeData(feeData.gasPrice, feeData.maxFeePerGas, feeData.maxPriorityFeePerGas);
      }
      return await getFeeDataFunc();
    }
    async estimateGas(_tx) {
      let tx = this._getTransactionRequest(_tx);
      if (isPromise(tx)) {
        tx = await tx;
      }
      return getBigInt(await this.#perform({
        method: "estimateGas",
        transaction: tx
      }), "%response");
    }
    async #call(tx, blockTag, attempt) {
      assert(attempt < MAX_CCIP_REDIRECTS, "CCIP read exceeded maximum redirections", "OFFCHAIN_FAULT", {
        reason: "TOO_MANY_REDIRECTS",
        transaction: Object.assign({}, tx, { blockTag, enableCcipRead: true })
      });
      const transaction = copyRequest(tx);
      try {
        return hexlify(await this._perform({ method: "call", transaction, blockTag }));
      } catch (error) {
        if (!this.disableCcipRead && isCallException(error) && error.data && attempt >= 0 && blockTag === "latest" && transaction.to != null && dataSlice(error.data, 0, 4) === "0x556f1830") {
          const data = error.data;
          const txSender = await resolveAddress(transaction.to, this);
          let ccipArgs;
          try {
            ccipArgs = parseOffchainLookup(dataSlice(error.data, 4));
          } catch (error2) {
            assert(false, error2.message, "OFFCHAIN_FAULT", {
              reason: "BAD_DATA",
              transaction,
              info: { data }
            });
          }
          assert(ccipArgs.sender.toLowerCase() === txSender.toLowerCase(), "CCIP Read sender mismatch", "CALL_EXCEPTION", {
            action: "call",
            data,
            reason: "OffchainLookup",
            transaction,
            invocation: null,
            revert: {
              signature: "OffchainLookup(address,string[],bytes,bytes4,bytes)",
              name: "OffchainLookup",
              args: ccipArgs.errorArgs
            }
          });
          const ccipResult = await this.ccipReadFetch(transaction, ccipArgs.calldata, ccipArgs.urls);
          assert(ccipResult != null, "CCIP Read failed to fetch data", "OFFCHAIN_FAULT", {
            reason: "FETCH_FAILED",
            transaction,
            info: { data: error.data, errorArgs: ccipArgs.errorArgs }
          });
          const tx2 = {
            to: txSender,
            data: concat([ccipArgs.selector, encodeBytes([ccipResult, ccipArgs.extraData])])
          };
          this.emit("debug", { action: "sendCcipReadCall", transaction: tx2 });
          try {
            const result = await this.#call(tx2, blockTag, attempt + 1);
            this.emit("debug", { action: "receiveCcipReadCallResult", transaction: Object.assign({}, tx2), result });
            return result;
          } catch (error2) {
            this.emit("debug", { action: "receiveCcipReadCallError", transaction: Object.assign({}, tx2), error: error2 });
            throw error2;
          }
        }
        throw error;
      }
    }
    async #checkNetwork(promise) {
      const { value } = await resolveProperties({
        network: this.getNetwork(),
        value: promise
      });
      return value;
    }
    async call(_tx) {
      const { tx, blockTag } = await resolveProperties({
        tx: this._getTransactionRequest(_tx),
        blockTag: this._getBlockTag(_tx.blockTag)
      });
      return await this.#checkNetwork(this.#call(tx, blockTag, _tx.enableCcipRead ? 0 : -1));
    }
    // Account
    async #getAccountValue(request, _address, _blockTag) {
      let address = this._getAddress(_address);
      let blockTag = this._getBlockTag(_blockTag);
      if (typeof address !== "string" || typeof blockTag !== "string") {
        [address, blockTag] = await Promise.all([address, blockTag]);
      }
      return await this.#checkNetwork(this.#perform(Object.assign(request, { address, blockTag })));
    }
    async getBalance(address, blockTag) {
      return getBigInt(await this.#getAccountValue({ method: "getBalance" }, address, blockTag), "%response");
    }
    async getTransactionCount(address, blockTag) {
      return getNumber(await this.#getAccountValue({ method: "getTransactionCount" }, address, blockTag), "%response");
    }
    async getCode(address, blockTag) {
      return hexlify(await this.#getAccountValue({ method: "getCode" }, address, blockTag));
    }
    async getStorage(address, _position, blockTag) {
      const position = getBigInt(_position, "position");
      return hexlify(await this.#getAccountValue({ method: "getStorage", position }, address, blockTag));
    }
    // Write
    async broadcastTransaction(signedTx) {
      const { blockNumber, hash: hash2, network } = await resolveProperties({
        blockNumber: this.getBlockNumber(),
        hash: this._perform({
          method: "broadcastTransaction",
          signedTransaction: signedTx
        }),
        network: this.getNetwork()
      });
      const tx = Transaction.from(signedTx);
      if (tx.hash !== hash2) {
        throw new Error("@TODO: the returned hash did not match");
      }
      return this._wrapTransactionResponse(tx, network).replaceableTransaction(blockNumber);
    }
    async #getBlock(block, includeTransactions) {
      if (isHexString(block, 32)) {
        return await this.#perform({
          method: "getBlock",
          blockHash: block,
          includeTransactions
        });
      }
      let blockTag = this._getBlockTag(block);
      if (typeof blockTag !== "string") {
        blockTag = await blockTag;
      }
      return await this.#perform({
        method: "getBlock",
        blockTag,
        includeTransactions
      });
    }
    // Queries
    async getBlock(block, prefetchTxs) {
      const { network, params } = await resolveProperties({
        network: this.getNetwork(),
        params: this.#getBlock(block, !!prefetchTxs)
      });
      if (params == null) {
        return null;
      }
      return this._wrapBlock(params, network);
    }
    async getTransaction(hash2) {
      const { network, params } = await resolveProperties({
        network: this.getNetwork(),
        params: this.#perform({ method: "getTransaction", hash: hash2 })
      });
      if (params == null) {
        return null;
      }
      return this._wrapTransactionResponse(params, network);
    }
    async getTransactionReceipt(hash2) {
      const { network, params } = await resolveProperties({
        network: this.getNetwork(),
        params: this.#perform({ method: "getTransactionReceipt", hash: hash2 })
      });
      if (params == null) {
        return null;
      }
      if (params.gasPrice == null && params.effectiveGasPrice == null) {
        const tx = await this.#perform({ method: "getTransaction", hash: hash2 });
        if (tx == null) {
          throw new Error("report this; could not find tx or effectiveGasPrice");
        }
        params.effectiveGasPrice = tx.gasPrice;
      }
      return this._wrapTransactionReceipt(params, network);
    }
    async getTransactionResult(hash2) {
      const { result } = await resolveProperties({
        network: this.getNetwork(),
        result: this.#perform({ method: "getTransactionResult", hash: hash2 })
      });
      if (result == null) {
        return null;
      }
      return hexlify(result);
    }
    // Bloom-filter Queries
    async getLogs(_filter) {
      let filter = this._getFilter(_filter);
      if (isPromise(filter)) {
        filter = await filter;
      }
      const { network, params } = await resolveProperties({
        network: this.getNetwork(),
        params: this.#perform({ method: "getLogs", filter })
      });
      return params.map((p4) => this._wrapLog(p4, network));
    }
    // ENS
    _getProvider(chainId) {
      assert(false, "provider cannot connect to target network", "UNSUPPORTED_OPERATION", {
        operation: "_getProvider()"
      });
    }
    async getResolver(name) {
      return await EnsResolver.fromName(this, name);
    }
    async getAvatar(name) {
      const resolver = await this.getResolver(name);
      if (resolver) {
        return await resolver.getAvatar();
      }
      return null;
    }
    async resolveName(name) {
      const resolver = await this.getResolver(name);
      if (resolver) {
        return await resolver.getAddress();
      }
      return null;
    }
    async lookupAddress(address) {
      address = getAddress(address);
      const node = namehash(address.substring(2).toLowerCase() + ".addr.reverse");
      try {
        const ensAddr = await EnsResolver.getEnsAddress(this);
        const ensContract = new Contract(ensAddr, [
          "function resolver(bytes32) view returns (address)"
        ], this);
        const resolver = await ensContract.resolver(node);
        if (resolver == null || resolver === ZeroAddress) {
          return null;
        }
        const resolverContract = new Contract(resolver, [
          "function name(bytes32) view returns (string)"
        ], this);
        const name = await resolverContract.name(node);
        const check = await this.resolveName(name);
        if (check !== address) {
          return null;
        }
        return name;
      } catch (error) {
        if (isError(error, "BAD_DATA") && error.value === "0x") {
          return null;
        }
        if (isError(error, "CALL_EXCEPTION")) {
          return null;
        }
        throw error;
      }
      return null;
    }
    async waitForTransaction(hash2, _confirms, timeout) {
      const confirms = _confirms != null ? _confirms : 1;
      if (confirms === 0) {
        return this.getTransactionReceipt(hash2);
      }
      return new Promise(async (resolve, reject) => {
        let timer = null;
        const listener = async (blockNumber) => {
          try {
            const receipt = await this.getTransactionReceipt(hash2);
            if (receipt != null) {
              if (blockNumber - receipt.blockNumber + 1 >= confirms) {
                resolve(receipt);
                if (timer) {
                  clearTimeout(timer);
                  timer = null;
                }
                return;
              }
            }
          } catch (error) {
            console.log("EEE", error);
          }
          this.once("block", listener);
        };
        if (timeout != null) {
          timer = setTimeout(() => {
            if (timer == null) {
              return;
            }
            timer = null;
            this.off("block", listener);
            reject(makeError("timeout", "TIMEOUT", { reason: "timeout" }));
          }, timeout);
        }
        listener(await this.getBlockNumber());
      });
    }
    async waitForBlock(blockTag) {
      assert(false, "not implemented yet", "NOT_IMPLEMENTED", {
        operation: "waitForBlock"
      });
    }
    /**
     *  Clear a timer created using the [[_setTimeout]] method.
     */
    _clearTimeout(timerId) {
      const timer = this.#timers.get(timerId);
      if (!timer) {
        return;
      }
      if (timer.timer) {
        clearTimeout(timer.timer);
      }
      this.#timers.delete(timerId);
    }
    /**
     *  Create a timer that will execute %%func%% after at least %%timeout%%
     *  (in ms). If %%timeout%% is unspecified, then %%func%% will execute
     *  in the next event loop.
     *
     *  [Pausing](AbstractProvider-paused) the provider will pause any
     *  associated timers.
     */
    _setTimeout(_func, timeout) {
      if (timeout == null) {
        timeout = 0;
      }
      const timerId = this.#nextTimer++;
      const func = () => {
        this.#timers.delete(timerId);
        _func();
      };
      if (this.paused) {
        this.#timers.set(timerId, { timer: null, func, time: timeout });
      } else {
        const timer = setTimeout(func, timeout);
        this.#timers.set(timerId, { timer, func, time: getTime2() });
      }
      return timerId;
    }
    /**
     *  Perform %%func%% on each subscriber.
     */
    _forEachSubscriber(func) {
      for (const sub of this.#subs.values()) {
        func(sub.subscriber);
      }
    }
    /**
     *  Sub-classes may override this to customize subscription
     *  implementations.
     */
    _getSubscriber(sub) {
      switch (sub.type) {
        case "debug":
        case "error":
        case "network":
          return new UnmanagedSubscriber(sub.type);
        case "block": {
          const subscriber = new PollingBlockSubscriber(this);
          subscriber.pollingInterval = this.pollingInterval;
          return subscriber;
        }
        case "event":
          return new PollingEventSubscriber(this, sub.filter);
        case "transaction":
          return new PollingTransactionSubscriber(this, sub.hash);
        case "orphan":
          return new PollingOrphanSubscriber(this, sub.filter);
      }
      throw new Error(`unsupported event: ${sub.type}`);
    }
    /**
     *  If a [[Subscriber]] fails and needs to replace itself, this
     *  method may be used.
     *
     *  For example, this is used for providers when using the
     *  ``eth_getFilterChanges`` method, which can return null if state
     *  filters are not supported by the backend, allowing the Subscriber
     *  to swap in a [[PollingEventSubscriber]].
     */
    _recoverSubscriber(oldSub, newSub) {
      for (const sub of this.#subs.values()) {
        if (sub.subscriber === oldSub) {
          if (sub.started) {
            sub.subscriber.stop();
          }
          sub.subscriber = newSub;
          if (sub.started) {
            newSub.start();
          }
          if (this.#pausedState != null) {
            newSub.pause(this.#pausedState);
          }
          break;
        }
      }
    }
    async #hasSub(event, emitArgs) {
      let sub = await getSubscription(event, this);
      if (sub.type === "event" && emitArgs && emitArgs.length > 0 && emitArgs[0].removed === true) {
        sub = await getSubscription({ orphan: "drop-log", log: emitArgs[0] }, this);
      }
      return this.#subs.get(sub.tag) || null;
    }
    async #getSub(event) {
      const subscription = await getSubscription(event, this);
      const tag = subscription.tag;
      let sub = this.#subs.get(tag);
      if (!sub) {
        const subscriber = this._getSubscriber(subscription);
        const addressableMap = /* @__PURE__ */ new WeakMap();
        const nameMap = /* @__PURE__ */ new Map();
        sub = { subscriber, tag, addressableMap, nameMap, started: false, listeners: [] };
        this.#subs.set(tag, sub);
      }
      return sub;
    }
    async on(event, listener) {
      const sub = await this.#getSub(event);
      sub.listeners.push({ listener, once: false });
      if (!sub.started) {
        sub.subscriber.start();
        sub.started = true;
        if (this.#pausedState != null) {
          sub.subscriber.pause(this.#pausedState);
        }
      }
      return this;
    }
    async once(event, listener) {
      const sub = await this.#getSub(event);
      sub.listeners.push({ listener, once: true });
      if (!sub.started) {
        sub.subscriber.start();
        sub.started = true;
        if (this.#pausedState != null) {
          sub.subscriber.pause(this.#pausedState);
        }
      }
      return this;
    }
    async emit(event, ...args) {
      const sub = await this.#hasSub(event, args);
      if (!sub || sub.listeners.length === 0) {
        return false;
      }
      ;
      const count = sub.listeners.length;
      sub.listeners = sub.listeners.filter(({ listener, once }) => {
        const payload = new EventPayload(this, once ? null : listener, event);
        try {
          listener.call(this, ...args, payload);
        } catch (error) {
        }
        return !once;
      });
      if (sub.listeners.length === 0) {
        if (sub.started) {
          sub.subscriber.stop();
        }
        this.#subs.delete(sub.tag);
      }
      return count > 0;
    }
    async listenerCount(event) {
      if (event) {
        const sub = await this.#hasSub(event);
        if (!sub) {
          return 0;
        }
        return sub.listeners.length;
      }
      let total = 0;
      for (const { listeners } of this.#subs.values()) {
        total += listeners.length;
      }
      return total;
    }
    async listeners(event) {
      if (event) {
        const sub = await this.#hasSub(event);
        if (!sub) {
          return [];
        }
        return sub.listeners.map(({ listener }) => listener);
      }
      let result = [];
      for (const { listeners } of this.#subs.values()) {
        result = result.concat(listeners.map(({ listener }) => listener));
      }
      return result;
    }
    async off(event, listener) {
      const sub = await this.#hasSub(event);
      if (!sub) {
        return this;
      }
      if (listener) {
        const index = sub.listeners.map(({ listener: listener2 }) => listener2).indexOf(listener);
        if (index >= 0) {
          sub.listeners.splice(index, 1);
        }
      }
      if (!listener || sub.listeners.length === 0) {
        if (sub.started) {
          sub.subscriber.stop();
        }
        this.#subs.delete(sub.tag);
      }
      return this;
    }
    async removeAllListeners(event) {
      if (event) {
        const { tag, started, subscriber } = await this.#getSub(event);
        if (started) {
          subscriber.stop();
        }
        this.#subs.delete(tag);
      } else {
        for (const [tag, { started, subscriber }] of this.#subs) {
          if (started) {
            subscriber.stop();
          }
          this.#subs.delete(tag);
        }
      }
      return this;
    }
    // Alias for "on"
    async addListener(event, listener) {
      return await this.on(event, listener);
    }
    // Alias for "off"
    async removeListener(event, listener) {
      return this.off(event, listener);
    }
    /**
     *  If this provider has been destroyed using the [[destroy]] method.
     *
     *  Once destroyed, all resources are reclaimed, internal event loops
     *  and timers are cleaned up and no further requests may be sent to
     *  the provider.
     */
    get destroyed() {
      return this.#destroyed;
    }
    /**
     *  Sub-classes may use this to shutdown any sockets or release their
     *  resources and reject any pending requests.
     *
     *  Sub-classes **must** call ``super.destroy()``.
     */
    destroy() {
      this.removeAllListeners();
      for (const timerId of this.#timers.keys()) {
        this._clearTimeout(timerId);
      }
      this.#destroyed = true;
    }
    /**
     *  Whether the provider is currently paused.
     *
     *  A paused provider will not emit any events, and generally should
     *  not make any requests to the network, but that is up to sub-classes
     *  to manage.
     *
     *  Setting ``paused = true`` is identical to calling ``.pause(false)``,
     *  which will buffer any events that occur while paused until the
     *  provider is unpaused.
     */
    get paused() {
      return this.#pausedState != null;
    }
    set paused(pause) {
      if (!!pause === this.paused) {
        return;
      }
      if (this.paused) {
        this.resume();
      } else {
        this.pause(false);
      }
    }
    /**
     *  Pause the provider. If %%dropWhilePaused%%, any events that occur
     *  while paused are dropped, otherwise all events will be emitted once
     *  the provider is unpaused.
     */
    pause(dropWhilePaused) {
      this.#lastBlockNumber = -1;
      if (this.#pausedState != null) {
        if (this.#pausedState == !!dropWhilePaused) {
          return;
        }
        assert(false, "cannot change pause type; resume first", "UNSUPPORTED_OPERATION", {
          operation: "pause"
        });
      }
      this._forEachSubscriber((s4) => s4.pause(dropWhilePaused));
      this.#pausedState = !!dropWhilePaused;
      for (const timer of this.#timers.values()) {
        if (timer.timer) {
          clearTimeout(timer.timer);
        }
        timer.time = getTime2() - timer.time;
      }
    }
    /**
     *  Resume the provider.
     */
    resume() {
      if (this.#pausedState == null) {
        return;
      }
      this._forEachSubscriber((s4) => s4.resume());
      this.#pausedState = null;
      for (const timer of this.#timers.values()) {
        let timeout = timer.time;
        if (timeout < 0) {
          timeout = 0;
        }
        timer.time = getTime2();
        setTimeout(timer.func, timeout);
      }
    }
  };
  function _parseString(result, start) {
    try {
      const bytes2 = _parseBytes(result, start);
      if (bytes2) {
        return toUtf8String(bytes2);
      }
    } catch (error) {
    }
    return null;
  }
  function _parseBytes(result, start) {
    if (result === "0x") {
      return null;
    }
    try {
      const offset = getNumber(dataSlice(result, start, start + 32));
      const length = getNumber(dataSlice(result, offset, offset + 32));
      return dataSlice(result, offset + 32, offset + 32 + length);
    } catch (error) {
    }
    return null;
  }
  function numPad(value) {
    const result = toBeArray(value);
    if (result.length > 32) {
      throw new Error("internal; should not happen");
    }
    const padded = new Uint8Array(32);
    padded.set(result, 32 - result.length);
    return padded;
  }
  function bytesPad(value) {
    if (value.length % 32 === 0) {
      return value;
    }
    const result = new Uint8Array(Math.ceil(value.length / 32) * 32);
    result.set(value);
    return result;
  }
  var empty = new Uint8Array([]);
  function encodeBytes(datas) {
    const result = [];
    let byteCount = 0;
    for (let i4 = 0; i4 < datas.length; i4++) {
      result.push(empty);
      byteCount += 32;
    }
    for (let i4 = 0; i4 < datas.length; i4++) {
      const data = getBytes(datas[i4]);
      result[i4] = numPad(byteCount);
      result.push(numPad(data.length));
      result.push(bytesPad(data));
      byteCount += 32 + Math.ceil(data.length / 32) * 32;
    }
    return concat(result);
  }
  var zeros = "0x0000000000000000000000000000000000000000000000000000000000000000";
  function parseOffchainLookup(data) {
    const result = {
      sender: "",
      urls: [],
      calldata: "",
      selector: "",
      extraData: "",
      errorArgs: []
    };
    assert(dataLength(data) >= 5 * 32, "insufficient OffchainLookup data", "OFFCHAIN_FAULT", {
      reason: "insufficient OffchainLookup data"
    });
    const sender = dataSlice(data, 0, 32);
    assert(dataSlice(sender, 0, 12) === dataSlice(zeros, 0, 12), "corrupt OffchainLookup sender", "OFFCHAIN_FAULT", {
      reason: "corrupt OffchainLookup sender"
    });
    result.sender = dataSlice(sender, 12);
    try {
      const urls = [];
      const urlsOffset = getNumber(dataSlice(data, 32, 64));
      const urlsLength = getNumber(dataSlice(data, urlsOffset, urlsOffset + 32));
      const urlsData = dataSlice(data, urlsOffset + 32);
      for (let u4 = 0; u4 < urlsLength; u4++) {
        const url = _parseString(urlsData, u4 * 32);
        if (url == null) {
          throw new Error("abort");
        }
        urls.push(url);
      }
      result.urls = urls;
    } catch (error) {
      assert(false, "corrupt OffchainLookup urls", "OFFCHAIN_FAULT", {
        reason: "corrupt OffchainLookup urls"
      });
    }
    try {
      const calldata = _parseBytes(data, 64);
      if (calldata == null) {
        throw new Error("abort");
      }
      result.calldata = calldata;
    } catch (error) {
      assert(false, "corrupt OffchainLookup calldata", "OFFCHAIN_FAULT", {
        reason: "corrupt OffchainLookup calldata"
      });
    }
    assert(dataSlice(data, 100, 128) === dataSlice(zeros, 0, 28), "corrupt OffchainLookup callbaackSelector", "OFFCHAIN_FAULT", {
      reason: "corrupt OffchainLookup callbaackSelector"
    });
    result.selector = dataSlice(data, 96, 100);
    try {
      const extraData = _parseBytes(data, 128);
      if (extraData == null) {
        throw new Error("abort");
      }
      result.extraData = extraData;
    } catch (error) {
      assert(false, "corrupt OffchainLookup extraData", "OFFCHAIN_FAULT", {
        reason: "corrupt OffchainLookup extraData"
      });
    }
    result.errorArgs = "sender,urls,calldata,selector,extraData".split(/,/).map((k3) => result[k3]);
    return result;
  }

  // node_modules/ethers/lib.esm/providers/abstract-signer.js
  function checkProvider(signer, operation) {
    if (signer.provider) {
      return signer.provider;
    }
    assert(false, "missing provider", "UNSUPPORTED_OPERATION", { operation });
  }
  async function populate(signer, tx) {
    let pop = copyRequest(tx);
    if (pop.to != null) {
      pop.to = resolveAddress(pop.to, signer);
    }
    if (pop.from != null) {
      const from = pop.from;
      pop.from = Promise.all([
        signer.getAddress(),
        resolveAddress(from, signer)
      ]).then(([address, from2]) => {
        assertArgument(address.toLowerCase() === from2.toLowerCase(), "transaction from mismatch", "tx.from", from2);
        return address;
      });
    } else {
      pop.from = signer.getAddress();
    }
    return await resolveProperties(pop);
  }
  var AbstractSigner = class {
    /**
     *  The provider this signer is connected to.
     */
    provider;
    /**
     *  Creates a new Signer connected to %%provider%%.
     */
    constructor(provider) {
      defineProperties(this, { provider: provider || null });
    }
    async getNonce(blockTag) {
      return checkProvider(this, "getTransactionCount").getTransactionCount(await this.getAddress(), blockTag);
    }
    async populateCall(tx) {
      const pop = await populate(this, tx);
      return pop;
    }
    async populateTransaction(tx) {
      const provider = checkProvider(this, "populateTransaction");
      const pop = await populate(this, tx);
      if (pop.nonce == null) {
        pop.nonce = await this.getNonce("pending");
      }
      if (pop.gasLimit == null) {
        pop.gasLimit = await this.estimateGas(pop);
      }
      const network = await this.provider.getNetwork();
      if (pop.chainId != null) {
        const chainId = getBigInt(pop.chainId);
        assertArgument(chainId === network.chainId, "transaction chainId mismatch", "tx.chainId", tx.chainId);
      } else {
        pop.chainId = network.chainId;
      }
      const hasEip1559 = pop.maxFeePerGas != null || pop.maxPriorityFeePerGas != null;
      if (pop.gasPrice != null && (pop.type === 2 || hasEip1559)) {
        assertArgument(false, "eip-1559 transaction do not support gasPrice", "tx", tx);
      } else if ((pop.type === 0 || pop.type === 1) && hasEip1559) {
        assertArgument(false, "pre-eip-1559 transaction do not support maxFeePerGas/maxPriorityFeePerGas", "tx", tx);
      }
      if ((pop.type === 2 || pop.type == null) && (pop.maxFeePerGas != null && pop.maxPriorityFeePerGas != null)) {
        pop.type = 2;
      } else if (pop.type === 0 || pop.type === 1) {
        const feeData = await provider.getFeeData();
        assert(feeData.gasPrice != null, "network does not support gasPrice", "UNSUPPORTED_OPERATION", {
          operation: "getGasPrice"
        });
        if (pop.gasPrice == null) {
          pop.gasPrice = feeData.gasPrice;
        }
      } else {
        const feeData = await provider.getFeeData();
        if (pop.type == null) {
          if (feeData.maxFeePerGas != null && feeData.maxPriorityFeePerGas != null) {
            pop.type = 2;
            if (pop.gasPrice != null) {
              const gasPrice = pop.gasPrice;
              delete pop.gasPrice;
              pop.maxFeePerGas = gasPrice;
              pop.maxPriorityFeePerGas = gasPrice;
            } else {
              if (pop.maxFeePerGas == null) {
                pop.maxFeePerGas = feeData.maxFeePerGas;
              }
              if (pop.maxPriorityFeePerGas == null) {
                pop.maxPriorityFeePerGas = feeData.maxPriorityFeePerGas;
              }
            }
          } else if (feeData.gasPrice != null) {
            assert(!hasEip1559, "network does not support EIP-1559", "UNSUPPORTED_OPERATION", {
              operation: "populateTransaction"
            });
            if (pop.gasPrice == null) {
              pop.gasPrice = feeData.gasPrice;
            }
            pop.type = 0;
          } else {
            assert(false, "failed to get consistent fee data", "UNSUPPORTED_OPERATION", {
              operation: "signer.getFeeData"
            });
          }
        } else if (pop.type === 2) {
          if (pop.maxFeePerGas == null) {
            pop.maxFeePerGas = feeData.maxFeePerGas;
          }
          if (pop.maxPriorityFeePerGas == null) {
            pop.maxPriorityFeePerGas = feeData.maxPriorityFeePerGas;
          }
        }
      }
      return await resolveProperties(pop);
    }
    async estimateGas(tx) {
      return checkProvider(this, "estimateGas").estimateGas(await this.populateCall(tx));
    }
    async call(tx) {
      return checkProvider(this, "call").call(await this.populateCall(tx));
    }
    async resolveName(name) {
      const provider = checkProvider(this, "resolveName");
      return await provider.resolveName(name);
    }
    async sendTransaction(tx) {
      const provider = checkProvider(this, "sendTransaction");
      const pop = await this.populateTransaction(tx);
      delete pop.from;
      const txObj = Transaction.from(pop);
      return await provider.broadcastTransaction(await this.signTransaction(txObj));
    }
  };

  // node_modules/ethers/lib.esm/providers/subscriber-filterid.js
  function copy3(obj) {
    return JSON.parse(JSON.stringify(obj));
  }
  var FilterIdSubscriber = class {
    #provider;
    #filterIdPromise;
    #poller;
    #running;
    #network;
    #hault;
    /**
     *  Creates a new **FilterIdSubscriber** which will used [[_subscribe]]
     *  and [[_emitResults]] to setup the subscription and provide the event
     *  to the %%provider%%.
     */
    constructor(provider) {
      this.#provider = provider;
      this.#filterIdPromise = null;
      this.#poller = this.#poll.bind(this);
      this.#running = false;
      this.#network = null;
      this.#hault = false;
    }
    /**
     *  Sub-classes **must** override this to begin the subscription.
     */
    _subscribe(provider) {
      throw new Error("subclasses must override this");
    }
    /**
     *  Sub-classes **must** override this handle the events.
     */
    _emitResults(provider, result) {
      throw new Error("subclasses must override this");
    }
    /**
     *  Sub-classes **must** override this handle recovery on errors.
     */
    _recover(provider) {
      throw new Error("subclasses must override this");
    }
    async #poll(blockNumber) {
      try {
        if (this.#filterIdPromise == null) {
          this.#filterIdPromise = this._subscribe(this.#provider);
        }
        let filterId = null;
        try {
          filterId = await this.#filterIdPromise;
        } catch (error) {
          if (!isError(error, "UNSUPPORTED_OPERATION") || error.operation !== "eth_newFilter") {
            throw error;
          }
        }
        if (filterId == null) {
          this.#filterIdPromise = null;
          this.#provider._recoverSubscriber(this, this._recover(this.#provider));
          return;
        }
        const network = await this.#provider.getNetwork();
        if (!this.#network) {
          this.#network = network;
        }
        if (this.#network.chainId !== network.chainId) {
          throw new Error("chaid changed");
        }
        if (this.#hault) {
          return;
        }
        const result = await this.#provider.send("eth_getFilterChanges", [filterId]);
        await this._emitResults(this.#provider, result);
      } catch (error) {
        console.log("@TODO", error);
      }
      this.#provider.once("block", this.#poller);
    }
    #teardown() {
      const filterIdPromise = this.#filterIdPromise;
      if (filterIdPromise) {
        this.#filterIdPromise = null;
        filterIdPromise.then((filterId) => {
          this.#provider.send("eth_uninstallFilter", [filterId]);
        });
      }
    }
    start() {
      if (this.#running) {
        return;
      }
      this.#running = true;
      this.#poll(-2);
    }
    stop() {
      if (!this.#running) {
        return;
      }
      this.#running = false;
      this.#hault = true;
      this.#teardown();
      this.#provider.off("block", this.#poller);
    }
    pause(dropWhilePaused) {
      if (dropWhilePaused) {
        this.#teardown();
      }
      this.#provider.off("block", this.#poller);
    }
    resume() {
      this.start();
    }
  };
  var FilterIdEventSubscriber = class extends FilterIdSubscriber {
    #event;
    /**
     *  Creates a new **FilterIdEventSubscriber** attached to %%provider%%
     *  listening for %%filter%%.
     */
    constructor(provider, filter) {
      super(provider);
      this.#event = copy3(filter);
    }
    _recover(provider) {
      return new PollingEventSubscriber(provider, this.#event);
    }
    async _subscribe(provider) {
      const filterId = await provider.send("eth_newFilter", [this.#event]);
      return filterId;
    }
    async _emitResults(provider, results) {
      for (const result of results) {
        provider.emit(this.#event, provider._wrapLog(result, provider._network));
      }
    }
  };
  var FilterIdPendingSubscriber = class extends FilterIdSubscriber {
    async _subscribe(provider) {
      return await provider.send("eth_newPendingTransactionFilter", []);
    }
    async _emitResults(provider, results) {
      for (const result of results) {
        provider.emit("pending", result);
      }
    }
  };

  // node_modules/ethers/lib.esm/providers/provider-jsonrpc.js
  var Primitive = "bigint,boolean,function,number,string,symbol".split(/,/g);
  function deepCopy(value) {
    if (value == null || Primitive.indexOf(typeof value) >= 0) {
      return value;
    }
    if (typeof value.getAddress === "function") {
      return value;
    }
    if (Array.isArray(value)) {
      return value.map(deepCopy);
    }
    if (typeof value === "object") {
      return Object.keys(value).reduce((accum, key) => {
        accum[key] = value[key];
        return accum;
      }, {});
    }
    throw new Error(`should not happen: ${value} (${typeof value})`);
  }
  function stall(duration) {
    return new Promise((resolve) => {
      setTimeout(resolve, duration);
    });
  }
  function getLowerCase(value) {
    if (value) {
      return value.toLowerCase();
    }
    return value;
  }
  function isPollable(value) {
    return value && typeof value.pollingInterval === "number";
  }
  var defaultOptions2 = {
    polling: false,
    staticNetwork: null,
    batchStallTime: 10,
    batchMaxSize: 1 << 20,
    batchMaxCount: 100,
    cacheTimeout: 250,
    pollingInterval: 4e3
  };
  var JsonRpcSigner = class extends AbstractSigner {
    address;
    constructor(provider, address) {
      super(provider);
      address = getAddress(address);
      defineProperties(this, { address });
    }
    connect(provider) {
      assert(false, "cannot reconnect JsonRpcSigner", "UNSUPPORTED_OPERATION", {
        operation: "signer.connect"
      });
    }
    async getAddress() {
      return this.address;
    }
    // JSON-RPC will automatially fill in nonce, etc. so we just check from
    async populateTransaction(tx) {
      return await this.populateCall(tx);
    }
    // Returns just the hash of the transaction after sent, which is what
    // the bare JSON-RPC API does;
    async sendUncheckedTransaction(_tx) {
      const tx = deepCopy(_tx);
      const promises = [];
      if (tx.from) {
        const _from = tx.from;
        promises.push((async () => {
          const from = await resolveAddress(_from, this.provider);
          assertArgument(from != null && from.toLowerCase() === this.address.toLowerCase(), "from address mismatch", "transaction", _tx);
          tx.from = from;
        })());
      } else {
        tx.from = this.address;
      }
      if (tx.gasLimit == null) {
        promises.push((async () => {
          tx.gasLimit = await this.provider.estimateGas({ ...tx, from: this.address });
        })());
      }
      if (tx.to != null) {
        const _to = tx.to;
        promises.push((async () => {
          tx.to = await resolveAddress(_to, this.provider);
        })());
      }
      if (promises.length) {
        await Promise.all(promises);
      }
      const hexTx = this.provider.getRpcTransaction(tx);
      return this.provider.send("eth_sendTransaction", [hexTx]);
    }
    async sendTransaction(tx) {
      const blockNumber = await this.provider.getBlockNumber();
      const hash2 = await this.sendUncheckedTransaction(tx);
      return await new Promise((resolve, reject) => {
        const timeouts = [1e3, 100];
        const checkTx = async () => {
          const tx2 = await this.provider.getTransaction(hash2);
          if (tx2 != null) {
            resolve(tx2.replaceableTransaction(blockNumber));
            return;
          }
          this.provider._setTimeout(() => {
            checkTx();
          }, timeouts.pop() || 4e3);
        };
        checkTx();
      });
    }
    async signTransaction(_tx) {
      const tx = deepCopy(_tx);
      if (tx.from) {
        const from = await resolveAddress(tx.from, this.provider);
        assertArgument(from != null && from.toLowerCase() === this.address.toLowerCase(), "from address mismatch", "transaction", _tx);
        tx.from = from;
      } else {
        tx.from = this.address;
      }
      const hexTx = this.provider.getRpcTransaction(tx);
      return await this.provider.send("eth_signTransaction", [hexTx]);
    }
    async signMessage(_message) {
      const message = typeof _message === "string" ? toUtf8Bytes(_message) : _message;
      return await this.provider.send("personal_sign", [
        hexlify(message),
        this.address.toLowerCase()
      ]);
    }
    async signTypedData(domain, types, _value) {
      const value = deepCopy(_value);
      const populated = await TypedDataEncoder.resolveNames(domain, types, value, async (value2) => {
        const address = await resolveAddress(value2);
        assertArgument(address != null, "TypedData does not support null address", "value", value2);
        return address;
      });
      return await this.provider.send("eth_signTypedData_v4", [
        this.address.toLowerCase(),
        JSON.stringify(TypedDataEncoder.getPayload(populated.domain, types, populated.value))
      ]);
    }
    async unlock(password) {
      return this.provider.send("personal_unlockAccount", [
        this.address.toLowerCase(),
        password,
        null
      ]);
    }
    // https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
    async _legacySignMessage(_message) {
      const message = typeof _message === "string" ? toUtf8Bytes(_message) : _message;
      return await this.provider.send("eth_sign", [
        this.address.toLowerCase(),
        hexlify(message)
      ]);
    }
  };
  var JsonRpcApiProvider = class extends AbstractProvider {
    #options;
    // The next ID to use for the JSON-RPC ID field
    #nextId;
    // Payloads are queued and triggered in batches using the drainTimer
    #payloads;
    #drainTimer;
    #notReady;
    #network;
    #scheduleDrain() {
      if (this.#drainTimer) {
        return;
      }
      const stallTime = this._getOption("batchMaxCount") === 1 ? 0 : this._getOption("batchStallTime");
      this.#drainTimer = setTimeout(() => {
        this.#drainTimer = null;
        const payloads = this.#payloads;
        this.#payloads = [];
        while (payloads.length) {
          const batch = [payloads.shift()];
          while (payloads.length) {
            if (batch.length === this.#options.batchMaxCount) {
              break;
            }
            batch.push(payloads.shift());
            const bytes2 = JSON.stringify(batch.map((p4) => p4.payload));
            if (bytes2.length > this.#options.batchMaxSize) {
              payloads.unshift(batch.pop());
              break;
            }
          }
          (async () => {
            const payload = batch.length === 1 ? batch[0].payload : batch.map((p4) => p4.payload);
            this.emit("debug", { action: "sendRpcPayload", payload });
            try {
              const result = await this._send(payload);
              this.emit("debug", { action: "receiveRpcResult", result });
              for (const { resolve, reject, payload: payload2 } of batch) {
                if (this.destroyed) {
                  reject(makeError("provider destroyed; cancelled request", "UNSUPPORTED_OPERATION", { operation: payload2.method }));
                  continue;
                }
                const resp = result.filter((r4) => r4.id === payload2.id)[0];
                if (resp == null) {
                  const error = makeError("missing response for request", "BAD_DATA", {
                    value: result,
                    info: { payload: payload2 }
                  });
                  this.emit("error", error);
                  reject(error);
                  continue;
                }
                if ("error" in resp) {
                  reject(this.getRpcError(payload2, resp));
                  continue;
                }
                resolve(resp.result);
              }
            } catch (error) {
              this.emit("debug", { action: "receiveRpcError", error });
              for (const { reject } of batch) {
                reject(error);
              }
            }
          })();
        }
      }, stallTime);
    }
    constructor(network, options) {
      super(network, options);
      this.#nextId = 1;
      this.#options = Object.assign({}, defaultOptions2, options || {});
      this.#payloads = [];
      this.#drainTimer = null;
      this.#network = null;
      {
        let resolve = null;
        const promise = new Promise((_resolve) => {
          resolve = _resolve;
        });
        this.#notReady = { promise, resolve };
      }
      const staticNetwork = this._getOption("staticNetwork");
      if (staticNetwork) {
        assertArgument(network == null || staticNetwork.matches(network), "staticNetwork MUST match network object", "options", options);
        this.#network = staticNetwork;
      }
    }
    /**
     *  Returns the value associated with the option %%key%%.
     *
     *  Sub-classes can use this to inquire about configuration options.
     */
    _getOption(key) {
      return this.#options[key];
    }
    /**
     *  Gets the [[Network]] this provider has committed to. On each call, the network
     *  is detected, and if it has changed, the call will reject.
     */
    get _network() {
      assert(this.#network, "network is not available yet", "NETWORK_ERROR");
      return this.#network;
    }
    /**
     *  Resolves to the non-normalized value by performing %%req%%.
     *
     *  Sub-classes may override this to modify behavior of actions,
     *  and should generally call ``super._perform`` as a fallback.
     */
    async _perform(req) {
      if (req.method === "call" || req.method === "estimateGas") {
        let tx = req.transaction;
        if (tx && tx.type != null && getBigInt(tx.type)) {
          if (tx.maxFeePerGas == null && tx.maxPriorityFeePerGas == null) {
            const feeData = await this.getFeeData();
            if (feeData.maxFeePerGas == null && feeData.maxPriorityFeePerGas == null) {
              req = Object.assign({}, req, {
                transaction: Object.assign({}, tx, { type: void 0 })
              });
            }
          }
        }
      }
      const request = this.getRpcRequest(req);
      if (request != null) {
        return await this.send(request.method, request.args);
      }
      return super._perform(req);
    }
    /**
     *  Sub-classes may override this; it detects the *actual* network that
     *  we are **currently** connected to.
     *
     *  Keep in mind that [[send]] may only be used once [[ready]], otherwise the
     *  _send primitive must be used instead.
     */
    async _detectNetwork() {
      const network = this._getOption("staticNetwork");
      if (network) {
        return network;
      }
      if (this.ready) {
        return Network.from(getBigInt(await this.send("eth_chainId", [])));
      }
      const payload = {
        id: this.#nextId++,
        method: "eth_chainId",
        params: [],
        jsonrpc: "2.0"
      };
      this.emit("debug", { action: "sendRpcPayload", payload });
      let result;
      try {
        result = (await this._send(payload))[0];
      } catch (error) {
        this.emit("debug", { action: "receiveRpcError", error });
        throw error;
      }
      this.emit("debug", { action: "receiveRpcResult", result });
      if ("result" in result) {
        return Network.from(getBigInt(result.result));
      }
      throw this.getRpcError(payload, result);
    }
    /**
     *  Sub-classes **MUST** call this. Until [[_start]] has been called, no calls
     *  will be passed to [[_send]] from [[send]]. If it is overridden, then
     *  ``super._start()`` **MUST** be called.
     *
     *  Calling it multiple times is safe and has no effect.
     */
    _start() {
      if (this.#notReady == null || this.#notReady.resolve == null) {
        return;
      }
      this.#notReady.resolve();
      this.#notReady = null;
      (async () => {
        while (this.#network == null && !this.destroyed) {
          try {
            this.#network = await this._detectNetwork();
          } catch (error) {
            if (this.destroyed) {
              break;
            }
            console.log("JsonRpcProvider failed to detect network and cannot start up; retry in 1s (perhaps the URL is wrong or the node is not started)");
            this.emit("error", makeError("failed to bootstrap network detection", "NETWORK_ERROR", { event: "initial-network-discovery", info: { error } }));
            await stall(1e3);
          }
        }
        this.#scheduleDrain();
      })();
    }
    /**
     *  Resolves once the [[_start]] has been called. This can be used in
     *  sub-classes to defer sending data until the connection has been
     *  established.
     */
    async _waitUntilReady() {
      if (this.#notReady == null) {
        return;
      }
      return await this.#notReady.promise;
    }
    /**
     *  Return a Subscriber that will manage the %%sub%%.
     *
     *  Sub-classes may override this to modify the behavior of
     *  subscription management.
     */
    _getSubscriber(sub) {
      if (sub.type === "pending") {
        return new FilterIdPendingSubscriber(this);
      }
      if (sub.type === "event") {
        if (this._getOption("polling")) {
          return new PollingEventSubscriber(this, sub.filter);
        }
        return new FilterIdEventSubscriber(this, sub.filter);
      }
      if (sub.type === "orphan" && sub.filter.orphan === "drop-log") {
        return new UnmanagedSubscriber("orphan");
      }
      return super._getSubscriber(sub);
    }
    /**
     *  Returns true only if the [[_start]] has been called.
     */
    get ready() {
      return this.#notReady == null;
    }
    /**
     *  Returns %%tx%% as a normalized JSON-RPC transaction request,
     *  which has all values hexlified and any numeric values converted
     *  to Quantity values.
     */
    getRpcTransaction(tx) {
      const result = {};
      ["chainId", "gasLimit", "gasPrice", "type", "maxFeePerGas", "maxPriorityFeePerGas", "nonce", "value"].forEach((key) => {
        if (tx[key] == null) {
          return;
        }
        let dstKey = key;
        if (key === "gasLimit") {
          dstKey = "gas";
        }
        result[dstKey] = toQuantity(getBigInt(tx[key], `tx.${key}`));
      });
      ["from", "to", "data"].forEach((key) => {
        if (tx[key] == null) {
          return;
        }
        result[key] = hexlify(tx[key]);
      });
      if (tx.accessList) {
        result["accessList"] = accessListify(tx.accessList);
      }
      return result;
    }
    /**
     *  Returns the request method and arguments required to perform
     *  %%req%%.
     */
    getRpcRequest(req) {
      switch (req.method) {
        case "chainId":
          return { method: "eth_chainId", args: [] };
        case "getBlockNumber":
          return { method: "eth_blockNumber", args: [] };
        case "getGasPrice":
          return { method: "eth_gasPrice", args: [] };
        case "getBalance":
          return {
            method: "eth_getBalance",
            args: [getLowerCase(req.address), req.blockTag]
          };
        case "getTransactionCount":
          return {
            method: "eth_getTransactionCount",
            args: [getLowerCase(req.address), req.blockTag]
          };
        case "getCode":
          return {
            method: "eth_getCode",
            args: [getLowerCase(req.address), req.blockTag]
          };
        case "getStorage":
          return {
            method: "eth_getStorageAt",
            args: [
              getLowerCase(req.address),
              "0x" + req.position.toString(16),
              req.blockTag
            ]
          };
        case "broadcastTransaction":
          return {
            method: "eth_sendRawTransaction",
            args: [req.signedTransaction]
          };
        case "getBlock":
          if ("blockTag" in req) {
            return {
              method: "eth_getBlockByNumber",
              args: [req.blockTag, !!req.includeTransactions]
            };
          } else if ("blockHash" in req) {
            return {
              method: "eth_getBlockByHash",
              args: [req.blockHash, !!req.includeTransactions]
            };
          }
          break;
        case "getTransaction":
          return {
            method: "eth_getTransactionByHash",
            args: [req.hash]
          };
        case "getTransactionReceipt":
          return {
            method: "eth_getTransactionReceipt",
            args: [req.hash]
          };
        case "call":
          return {
            method: "eth_call",
            args: [this.getRpcTransaction(req.transaction), req.blockTag]
          };
        case "estimateGas": {
          return {
            method: "eth_estimateGas",
            args: [this.getRpcTransaction(req.transaction)]
          };
        }
        case "getLogs":
          if (req.filter && req.filter.address != null) {
            if (Array.isArray(req.filter.address)) {
              req.filter.address = req.filter.address.map(getLowerCase);
            } else {
              req.filter.address = getLowerCase(req.filter.address);
            }
          }
          return { method: "eth_getLogs", args: [req.filter] };
      }
      return null;
    }
    /**
     *  Returns an ethers-style Error for the given JSON-RPC error
     *  %%payload%%, coalescing the various strings and error shapes
     *  that different nodes return, coercing them into a machine-readable
     *  standardized error.
     */
    getRpcError(payload, _error) {
      const { method } = payload;
      const { error } = _error;
      if (method === "eth_estimateGas" && error.message) {
        const msg = error.message;
        if (!msg.match(/revert/i) && msg.match(/insufficient funds/i)) {
          return makeError("insufficient funds", "INSUFFICIENT_FUNDS", {
            transaction: payload.params[0],
            info: { payload, error }
          });
        }
      }
      if (method === "eth_call" || method === "eth_estimateGas") {
        const result = spelunkData(error);
        const e4 = AbiCoder.getBuiltinCallException(method === "eth_call" ? "call" : "estimateGas", payload.params[0], result ? result.data : null);
        e4.info = { error, payload };
        return e4;
      }
      const message = JSON.stringify(spelunkMessage(error));
      if (typeof error.message === "string" && error.message.match(/user denied|ethers-user-denied/i)) {
        const actionMap = {
          eth_sign: "signMessage",
          personal_sign: "signMessage",
          eth_signTypedData_v4: "signTypedData",
          eth_signTransaction: "signTransaction",
          eth_sendTransaction: "sendTransaction",
          eth_requestAccounts: "requestAccess",
          wallet_requestAccounts: "requestAccess"
        };
        return makeError(`user rejected action`, "ACTION_REJECTED", {
          action: actionMap[method] || "unknown",
          reason: "rejected",
          info: { payload, error }
        });
      }
      if (method === "eth_sendRawTransaction" || method === "eth_sendTransaction") {
        const transaction = payload.params[0];
        if (message.match(/insufficient funds|base fee exceeds gas limit/i)) {
          return makeError("insufficient funds for intrinsic transaction cost", "INSUFFICIENT_FUNDS", {
            transaction,
            info: { error }
          });
        }
        if (message.match(/nonce/i) && message.match(/too low/i)) {
          return makeError("nonce has already been used", "NONCE_EXPIRED", { transaction, info: { error } });
        }
        if (message.match(/replacement transaction/i) && message.match(/underpriced/i)) {
          return makeError("replacement fee too low", "REPLACEMENT_UNDERPRICED", { transaction, info: { error } });
        }
        if (message.match(/only replay-protected/i)) {
          return makeError("legacy pre-eip-155 transactions not supported", "UNSUPPORTED_OPERATION", {
            operation: method,
            info: { transaction, info: { error } }
          });
        }
      }
      let unsupported = !!message.match(/the method .* does not exist/i);
      if (!unsupported) {
        if (error && error.details && error.details.startsWith("Unauthorized method:")) {
          unsupported = true;
        }
      }
      if (unsupported) {
        return makeError("unsupported operation", "UNSUPPORTED_OPERATION", {
          operation: payload.method,
          info: { error, payload }
        });
      }
      return makeError("could not coalesce error", "UNKNOWN_ERROR", { error, payload });
    }
    /**
     *  Requests the %%method%% with %%params%% via the JSON-RPC protocol
     *  over the underlying channel. This can be used to call methods
     *  on the backend that do not have a high-level API within the Provider
     *  API.
     *
     *  This method queues requests according to the batch constraints
     *  in the options, assigns the request a unique ID.
     *
     *  **Do NOT override** this method in sub-classes; instead
     *  override [[_send]] or force the options values in the
     *  call to the constructor to modify this method's behavior.
     */
    send(method, params) {
      if (this.destroyed) {
        return Promise.reject(makeError("provider destroyed; cancelled request", "UNSUPPORTED_OPERATION", { operation: method }));
      }
      const id2 = this.#nextId++;
      const promise = new Promise((resolve, reject) => {
        this.#payloads.push({
          resolve,
          reject,
          payload: { method, params, id: id2, jsonrpc: "2.0" }
        });
      });
      this.#scheduleDrain();
      return promise;
    }
    /**
     *  Resolves to the [[Signer]] account for  %%address%% managed by
     *  the client.
     *
     *  If the %%address%% is a number, it is used as an index in the
     *  the accounts from [[listAccounts]].
     *
     *  This can only be used on clients which manage accounts (such as
     *  Geth with imported account or MetaMask).
     *
     *  Throws if the account doesn't exist.
     */
    async getSigner(address) {
      if (address == null) {
        address = 0;
      }
      const accountsPromise = this.send("eth_accounts", []);
      if (typeof address === "number") {
        const accounts2 = await accountsPromise;
        if (address >= accounts2.length) {
          throw new Error("no such account");
        }
        return new JsonRpcSigner(this, accounts2[address]);
      }
      const { accounts } = await resolveProperties({
        network: this.getNetwork(),
        accounts: accountsPromise
      });
      address = getAddress(address);
      for (const account of accounts) {
        if (getAddress(account) === address) {
          return new JsonRpcSigner(this, address);
        }
      }
      throw new Error("invalid account");
    }
    async listAccounts() {
      const accounts = await this.send("eth_accounts", []);
      return accounts.map((a4) => new JsonRpcSigner(this, a4));
    }
    destroy() {
      if (this.#drainTimer) {
        clearTimeout(this.#drainTimer);
        this.#drainTimer = null;
      }
      for (const { payload, reject } of this.#payloads) {
        reject(makeError("provider destroyed; cancelled request", "UNSUPPORTED_OPERATION", { operation: payload.method }));
      }
      this.#payloads = [];
      super.destroy();
    }
  };
  var JsonRpcApiPollingProvider = class extends JsonRpcApiProvider {
    #pollingInterval;
    constructor(network, options) {
      super(network, options);
      this.#pollingInterval = 4e3;
    }
    _getSubscriber(sub) {
      const subscriber = super._getSubscriber(sub);
      if (isPollable(subscriber)) {
        subscriber.pollingInterval = this.#pollingInterval;
      }
      return subscriber;
    }
    /**
     *  The polling interval (default: 4000 ms)
     */
    get pollingInterval() {
      return this.#pollingInterval;
    }
    set pollingInterval(value) {
      if (!Number.isInteger(value) || value < 0) {
        throw new Error("invalid interval");
      }
      this.#pollingInterval = value;
      this._forEachSubscriber((sub) => {
        if (isPollable(sub)) {
          sub.pollingInterval = this.#pollingInterval;
        }
      });
    }
  };
  function spelunkData(value) {
    if (value == null) {
      return null;
    }
    if (typeof value.message === "string" && value.message.match(/revert/i) && isHexString(value.data)) {
      return { message: value.message, data: value.data };
    }
    if (typeof value === "object") {
      for (const key in value) {
        const result = spelunkData(value[key]);
        if (result) {
          return result;
        }
      }
      return null;
    }
    if (typeof value === "string") {
      try {
        return spelunkData(JSON.parse(value));
      } catch (error) {
      }
    }
    return null;
  }
  function _spelunkMessage(value, result) {
    if (value == null) {
      return;
    }
    if (typeof value.message === "string") {
      result.push(value.message);
    }
    if (typeof value === "object") {
      for (const key in value) {
        _spelunkMessage(value[key], result);
      }
    }
    if (typeof value === "string") {
      try {
        return _spelunkMessage(JSON.parse(value), result);
      } catch (error) {
      }
    }
  }
  function spelunkMessage(value) {
    const result = [];
    _spelunkMessage(value, result);
    return result;
  }

  // node_modules/ethers/lib.esm/providers/provider-browser.js
  var BrowserProvider = class extends JsonRpcApiPollingProvider {
    #request;
    /**
     *  Connnect to the %%ethereum%% provider, optionally forcing the
     *  %%network%%.
     */
    constructor(ethereum, network) {
      super(network, { batchMaxCount: 1 });
      this.#request = async (method, params) => {
        const payload = { method, params };
        this.emit("debug", { action: "sendEip1193Request", payload });
        try {
          const result = await ethereum.request(payload);
          this.emit("debug", { action: "receiveEip1193Result", result });
          return result;
        } catch (e4) {
          const error = new Error(e4.message);
          error.code = e4.code;
          error.data = e4.data;
          error.payload = payload;
          this.emit("debug", { action: "receiveEip1193Error", error });
          throw error;
        }
      };
    }
    async send(method, params) {
      await this._start();
      return await super.send(method, params);
    }
    async _send(payload) {
      assertArgument(!Array.isArray(payload), "EIP-1193 does not support batch request", "payload", payload);
      try {
        const result = await this.#request(payload.method, payload.params || []);
        return [{ id: payload.id, result }];
      } catch (e4) {
        return [{
          id: payload.id,
          error: { code: e4.code, data: e4.data, message: e4.message }
        }];
      }
    }
    getRpcError(payload, error) {
      error = JSON.parse(JSON.stringify(error));
      switch (error.error.code || -1) {
        case 4001:
          error.error.message = `ethers-user-denied: ${error.error.message}`;
          break;
        case 4200:
          error.error.message = `ethers-unsupported: ${error.error.message}`;
          break;
      }
      return super.getRpcError(payload, error);
    }
    /**
     *  Resolves to ``true`` if the provider manages the %%address%%.
     */
    async hasSigner(address) {
      if (address == null) {
        address = 0;
      }
      const accounts = await this.send("eth_accounts", []);
      if (typeof address === "number") {
        return accounts.length > address;
      }
      address = address.toLowerCase();
      return accounts.filter((a4) => a4.toLowerCase() === address).length !== 0;
    }
    async getSigner(address) {
      if (address == null) {
        address = 0;
      }
      if (!await this.hasSigner(address)) {
        try {
          await this.#request("eth_requestAccounts", []);
        } catch (error) {
          const payload = error.payload;
          throw this.getRpcError(payload, { id: payload.id, error });
        }
      }
      return await super.getSigner(address);
    }
  };

  // node_modules/goober/dist/goober.modern.js
  var e3 = { data: "" };
  var t3 = (t4) => "object" == typeof window ? ((t4 ? t4.querySelector("#_goober") : window._goober) || Object.assign((t4 || document.head).appendChild(document.createElement("style")), { innerHTML: " ", id: "_goober" })).firstChild : t4 || e3;
  var l3 = /(?:([\u0080-\uFFFF\w-%@]+) *:? *([^{;]+?);|([^;}{]*?) *{)|(}\s*)/g;
  var a3 = /\/\*[^]*?\*\/|  +/g;
  var n3 = /\n+/g;
  var o3 = (e4, t4) => {
    let r4 = "", l4 = "", a4 = "";
    for (let n4 in e4) {
      let c4 = e4[n4];
      "@" == n4[0] ? "i" == n4[1] ? r4 = n4 + " " + c4 + ";" : l4 += "f" == n4[1] ? o3(c4, n4) : n4 + "{" + o3(c4, "k" == n4[1] ? "" : t4) + "}" : "object" == typeof c4 ? l4 += o3(c4, t4 ? t4.replace(/([^,])+/g, (e5) => n4.replace(/(^:.*)|([^,])+/g, (t5) => /&/.test(t5) ? t5.replace(/&/g, e5) : e5 ? e5 + " " + t5 : t5)) : n4) : null != c4 && (n4 = /^--/.test(n4) ? n4 : n4.replace(/[A-Z]/g, "-$&").toLowerCase(), a4 += o3.p ? o3.p(n4, c4) : n4 + ":" + c4 + ";");
    }
    return r4 + (t4 && a4 ? t4 + "{" + a4 + "}" : a4) + l4;
  };
  var c3 = {};
  var s3 = (e4) => {
    if ("object" == typeof e4) {
      let t4 = "";
      for (let r4 in e4)
        t4 += r4 + s3(e4[r4]);
      return t4;
    }
    return e4;
  };
  var i3 = (e4, t4, r4, i4, p4) => {
    let u4 = s3(e4), d4 = c3[u4] || (c3[u4] = ((e5) => {
      let t5 = 0, r5 = 11;
      for (; t5 < e5.length; )
        r5 = 101 * r5 + e5.charCodeAt(t5++) >>> 0;
      return "go" + r5;
    })(u4));
    if (!c3[d4]) {
      let t5 = u4 !== e4 ? e4 : ((e5) => {
        let t6, r5, o4 = [{}];
        for (; t6 = l3.exec(e5.replace(a3, "")); )
          t6[4] ? o4.shift() : t6[3] ? (r5 = t6[3].replace(n3, " ").trim(), o4.unshift(o4[0][r5] = o4[0][r5] || {})) : o4[0][t6[1]] = t6[2].replace(n3, " ").trim();
        return o4[0];
      })(e4);
      c3[d4] = o3(p4 ? { ["@keyframes " + d4]: t5 } : t5, r4 ? "" : "." + d4);
    }
    let f4 = r4 && c3.g ? c3.g : null;
    return r4 && (c3.g = c3[d4]), ((e5, t5, r5, l4) => {
      l4 ? t5.data = t5.data.replace(l4, e5) : -1 === t5.data.indexOf(e5) && (t5.data = r5 ? e5 + t5.data : t5.data + e5);
    })(c3[d4], t4, i4, f4), d4;
  };
  var p3 = (e4, t4, r4) => e4.reduce((e5, l4, a4) => {
    let n4 = t4[a4];
    if (n4 && n4.call) {
      let e6 = n4(r4), t5 = e6 && e6.props && e6.props.className || /^go/.test(e6) && e6;
      n4 = t5 ? "." + t5 : e6 && "object" == typeof e6 ? e6.props ? "" : o3(e6, "") : false === e6 ? "" : e6;
    }
    return e5 + l4 + (null == n4 ? "" : n4);
  }, "");
  function u3(e4) {
    let r4 = this || {}, l4 = e4.call ? e4(r4.p) : e4;
    return i3(l4.unshift ? l4.raw ? p3(l4, [].slice.call(arguments, 1), r4.p) : l4.reduce((e5, t4) => Object.assign(e5, t4 && t4.call ? t4(r4.p) : t4), {}) : l4, t3(r4.target), r4.g, r4.o, r4.k);
  }
  var d3;
  var f3;
  var g3;
  var b4 = u3.bind({ g: 1 });
  var h3 = u3.bind({ k: 1 });
  function m3(e4, t4, r4, l4) {
    o3.p = t4, d3 = e4, f3 = r4, g3 = l4;
  }
  function j3(e4, t4) {
    let r4 = this || {};
    return function() {
      let l4 = arguments;
      function a4(n4, o4) {
        let c4 = Object.assign({}, n4), s4 = c4.className || a4.className;
        r4.p = Object.assign({ theme: f3 && f3() }, c4), r4.o = / *go\d+/.test(s4), c4.className = u3.apply(r4, l4) + (s4 ? " " + s4 : ""), t4 && (c4.ref = o4);
        let i4 = e4;
        return e4[0] && (i4 = c4.as || e4, delete c4.as), g3 && i4[0] && g3(c4), d3(i4, c4);
      }
      return t4 ? t4(a4) : a4;
    };
  }

  // src/index.js
  var import_detect_provider = __toESM(require_dist());

  // node_modules/js-base64/base64.mjs
  var _hasatob = typeof atob === "function";
  var _hasBuffer = typeof Buffer === "function";
  var _TD = typeof TextDecoder === "function" ? new TextDecoder() : void 0;
  var _TE = typeof TextEncoder === "function" ? new TextEncoder() : void 0;
  var b64ch = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=";
  var b64chs = Array.prototype.slice.call(b64ch);
  var b64tab = ((a4) => {
    let tab = {};
    a4.forEach((c4, i4) => tab[c4] = i4);
    return tab;
  })(b64chs);
  var b64re = /^(?:[A-Za-z\d+\/]{4})*?(?:[A-Za-z\d+\/]{2}(?:==)?|[A-Za-z\d+\/]{3}=?)?$/;
  var _fromCC = String.fromCharCode.bind(String);
  var _U8Afrom = typeof Uint8Array.from === "function" ? Uint8Array.from.bind(Uint8Array) : (it) => new Uint8Array(Array.prototype.slice.call(it, 0));
  var _tidyB64 = (s4) => s4.replace(/[^A-Za-z0-9\+\/]/g, "");
  var re_btou = /[\xC0-\xDF][\x80-\xBF]|[\xE0-\xEF][\x80-\xBF]{2}|[\xF0-\xF7][\x80-\xBF]{3}/g;
  var cb_btou = (cccc) => {
    switch (cccc.length) {
      case 4:
        var cp = (7 & cccc.charCodeAt(0)) << 18 | (63 & cccc.charCodeAt(1)) << 12 | (63 & cccc.charCodeAt(2)) << 6 | 63 & cccc.charCodeAt(3), offset = cp - 65536;
        return _fromCC((offset >>> 10) + 55296) + _fromCC((offset & 1023) + 56320);
      case 3:
        return _fromCC((15 & cccc.charCodeAt(0)) << 12 | (63 & cccc.charCodeAt(1)) << 6 | 63 & cccc.charCodeAt(2));
      default:
        return _fromCC((31 & cccc.charCodeAt(0)) << 6 | 63 & cccc.charCodeAt(1));
    }
  };
  var btou = (b5) => b5.replace(re_btou, cb_btou);
  var atobPolyfill = (asc) => {
    asc = asc.replace(/\s+/g, "");
    if (!b64re.test(asc))
      throw new TypeError("malformed base64.");
    asc += "==".slice(2 - (asc.length & 3));
    let u24, bin = "", r1, r22;
    for (let i4 = 0; i4 < asc.length; ) {
      u24 = b64tab[asc.charAt(i4++)] << 18 | b64tab[asc.charAt(i4++)] << 12 | (r1 = b64tab[asc.charAt(i4++)]) << 6 | (r22 = b64tab[asc.charAt(i4++)]);
      bin += r1 === 64 ? _fromCC(u24 >> 16 & 255) : r22 === 64 ? _fromCC(u24 >> 16 & 255, u24 >> 8 & 255) : _fromCC(u24 >> 16 & 255, u24 >> 8 & 255, u24 & 255);
    }
    return bin;
  };
  var _atob = _hasatob ? (asc) => atob(_tidyB64(asc)) : _hasBuffer ? (asc) => Buffer.from(asc, "base64").toString("binary") : atobPolyfill;
  var _toUint8Array = _hasBuffer ? (a4) => _U8Afrom(Buffer.from(a4, "base64")) : (a4) => _U8Afrom(_atob(a4).split("").map((c4) => c4.charCodeAt(0)));
  var _decode2 = _hasBuffer ? (a4) => Buffer.from(a4, "base64").toString("utf8") : _TD ? (a4) => _TD.decode(_toUint8Array(a4)) : (a4) => btou(_atob(a4));
  var _unURI = (a4) => _tidyB64(a4.replace(/[-_]/g, (m0) => m0 == "-" ? "+" : "/"));
  var decode = (src) => _decode2(_unURI(src));

  // src/bridge.js
  var bridgeAddress = "0xb6C735bfF2B23f20e2603D4394FE3aF3e2B1EB69";
  var bridgeAbi = [
    {
      "inputs": [
        {
          "internalType": "bool",
          "name": "testnet",
          "type": "bool"
        },
        {
          "internalType": "uint256",
          "name": "fee",
          "type": "uint256"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "constructor"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "trusteeAddress",
          "type": "address"
        }
      ],
      "name": "BridgeAddTrustee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "fee",
          "type": "uint256"
        }
      ],
      "name": "BridgeChangeFee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "uint32",
          "name": "quorum",
          "type": "uint32"
        }
      ],
      "name": "BridgeChangeQuorum",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "bytes32",
          "name": "deroTxHash",
          "type": "bytes32"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "toAddress",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "BridgeFilled",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "sponsor",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "name",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "symbol",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint8",
          "name": "decimals",
          "type": "uint8"
        }
      ],
      "name": "BridgePairRegister",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "address",
          "name": "trusteeAddress",
          "type": "address"
        }
      ],
      "name": "BridgeRemoveTrustee",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        },
        {
          "indexed": true,
          "internalType": "address",
          "name": "fromAddr",
          "type": "address"
        },
        {
          "indexed": false,
          "internalType": "string",
          "name": "deroWallet",
          "type": "string"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        },
        {
          "indexed": false,
          "internalType": "uint256",
          "name": "feeAmount",
          "type": "uint256"
        }
      ],
      "name": "BridgeStarted",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": false,
          "internalType": "bool",
          "name": "bridgeOpen",
          "type": "bool"
        }
      ],
      "name": "BridgeState",
      "type": "event"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        },
        {
          "internalType": "string",
          "name": "deroWallet",
          "type": "string"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "bridgeETH2DERO",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "payable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "bridgeFee",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "bridgeOpen",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "deroTestnet",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "getSymbolListLength",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "initiatives",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "proposal",
          "type": "uint256"
        },
        {
          "internalType": "uint256",
          "name": "tally",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "quorum",
      "outputs": [
        {
          "internalType": "uint8",
          "name": "",
          "type": "uint8"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        }
      ],
      "name": "registerBridgePairToDERO",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "registeredERC20",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "name": "registeredSymbol",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "setBridgeClosed",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "symbolList",
      "outputs": [
        {
          "internalType": "string",
          "name": "",
          "type": "string"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "name": "trusteeList",
      "outputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "",
          "type": "address"
        }
      ],
      "name": "trustees",
      "outputs": [
        {
          "internalType": "uint256",
          "name": "",
          "type": "uint256"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "trustee",
          "type": "address"
        }
      ],
      "name": "voteAddTrustee",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "bytes32",
          "name": "deroTxHash",
          "type": "bytes32"
        },
        {
          "internalType": "address",
          "name": "erc20Addr",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "toAddress",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "voteFillDERO2ETHBridge",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "address",
          "name": "trustee",
          "type": "address"
        }
      ],
      "name": "voteRemoveTrustee",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint256",
          "name": "fee",
          "type": "uint256"
        }
      ],
      "name": "voteSetBridgeFee",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [],
      "name": "voteSetBridgeOpen",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "uint8",
          "name": "q",
          "type": "uint8"
        }
      ],
      "name": "voteSetQuorum",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "nonpayable",
      "type": "function"
    }
  ];
  var tokenAbi = [
    {
      "constant": true,
      "inputs": [],
      "name": "name",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_spender",
          "type": "address"
        },
        {
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "approve",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "totalSupply",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_from",
          "type": "address"
        },
        {
          "name": "_to",
          "type": "address"
        },
        {
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "transferFrom",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "decimals",
      "outputs": [
        {
          "name": "",
          "type": "uint8"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "_owner",
          "type": "address"
        }
      ],
      "name": "balanceOf",
      "outputs": [
        {
          "name": "balance",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [],
      "name": "symbol",
      "outputs": [
        {
          "name": "",
          "type": "string"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "constant": false,
      "inputs": [
        {
          "name": "_to",
          "type": "address"
        },
        {
          "name": "_value",
          "type": "uint256"
        }
      ],
      "name": "transfer",
      "outputs": [
        {
          "name": "",
          "type": "bool"
        }
      ],
      "payable": false,
      "stateMutability": "nonpayable",
      "type": "function"
    },
    {
      "constant": true,
      "inputs": [
        {
          "name": "_owner",
          "type": "address"
        },
        {
          "name": "_spender",
          "type": "address"
        }
      ],
      "name": "allowance",
      "outputs": [
        {
          "name": "",
          "type": "uint256"
        }
      ],
      "payable": false,
      "stateMutability": "view",
      "type": "function"
    },
    {
      "payable": true,
      "stateMutability": "payable",
      "type": "fallback"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "owner",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "spender",
          "type": "address"
        },
        {
          "indexed": false,
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Approval",
      "type": "event"
    },
    {
      "anonymous": false,
      "inputs": [
        {
          "indexed": true,
          "name": "from",
          "type": "address"
        },
        {
          "indexed": true,
          "name": "to",
          "type": "address"
        },
        {
          "indexed": false,
          "name": "value",
          "type": "uint256"
        }
      ],
      "name": "Transfer",
      "type": "event"
    }
  ];

  // src/style.js
  b4`
  @font-face {
    font-family: GoNotoKurrent;
    src: url(./fonts/GoNotoKurrent-Regular.ttf) format('truetype');
  }

  html, body {
    font-family: GoNotoKurrent;

    /* thanks to https://www.magicpattern.design/tools/css-backgrounds */
    background-color: #1a1a1a;
    background-image: radial-gradient(#2c2c2c 1.35px, #1a1a1a 1.35px);
    background-size: 27px 27px;
    margin: 0 1em;
  }
`;
  var rotate = h3`
    from, to {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
`;
  var RotateIcon = j3("span")`
  animation: ${rotate} 1s ease-in-out infinite;
`;
  var Box = j3("div")`
  max-width: 350px;
  background-color: white;
  border-radius: 10px;
  padding: 2em;
  margin: 5em auto;
  position: relative;
`;
  var BoxTitle = j3("div")`
  font-size: 34px;
  font-weight: bold;
  text-align: center;
  margin: 10px 0;
`;
  var BoxDescription = j3("div")`
  font-size: 20px;
  text-align: center;
  margin: 10px 0 20px 0;
`;
  var BoxSubmitButton = j3("button")`
  border-radius: 10px;
  padding: 0.7em;
  font-size: 20px;
  font-weight: bold;
  border: none;
  width: 100%;
  cursor: pointer;
  margin-top: 20px;
  background-color: #f0f0f0;
  display: flex;
  align-items: center;
  gap: .5em;
  justify-content: center;
`;
  var ErrDiv = j3("div")`
  border: 3px solid red;
  border-radius: 10px;
  padding: 1em;
  color: gray;
  word-break: break-all;
`;
  var WarningDiv = j3("div")`
  border: 3px solid yellow;
  border-radius: 10px;
  padding: 1em;
  color: gray;
`;
  var SuccessDiv = j3("div")`
  border: 3px solid green;
  border-radius: 10px;
  padding: 1em;
  color: gray;
`;
  var BoxItemTitle = j3("div")`
  margin-bottom: 5px;
`;
  var BoxItemValue = j3("div")`
  margin-bottom: 10px;
  color: gray;
  word-break: break-all;
  font-size: 14px;
`;
  var BridgeTransfer = j3("div")`
  padding: 1em;
  border-radius: 10px;
  margin: 20px 0;
  font-size: 26px;
  border: 3px solid #f0f0f0;
  text-align: center;
  font-weight: bold;
`;
  var MetaMaskLogoWrap = j3("div")`
  padding: 15px;
  left: 50%;
  margin-left: -40px;
  top: -40px;
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background-color: white;
  position: absolute;
`;
  var MetaMaskLogo = j3("div")`
  background-image: url(./images/metamask.png);
  background-size: contain;
  background-repeat: no-repeat;
  width: 100%;
  height: 100%;
`;

  // src/index.js
  m3(y);
  var BoxItem = (props) => {
    const { title, value } = props;
    var render = value;
    if (typeof render === "function") {
      render = value();
    }
    return /* @__PURE__ */ y("div", null, /* @__PURE__ */ y(BoxItemTitle, null, title), /* @__PURE__ */ y(BoxItemValue, null, render));
  };
  var App = () => {
    const [loaded, setLoaded] = h2();
    const [err, setErr] = h2();
    const [chainId, setChainId] = h2();
    const [account, setAccount] = h2();
    const [bridgeData, setBridgeData] = h2();
    const [provider, setProvider] = h2();
    const [bridging, setBridging] = h2();
    const [bridgeFee, setBridgeFee] = h2();
    const [bridgeHash, setBridgeHash] = h2();
    p2(async () => {
      const provider2 = await (0, import_detect_provider.default)();
      if (!provider2) {
        setErr("Metamask is not available or multiple wallets installed?");
        return;
      }
      setProvider(new BrowserProvider(window.ethereum));
      try {
        const queryString = new URL(location).searchParams;
        const base64Data = decode(queryString.get("data"));
        const queryData = JSON.parse(base64Data);
        setBridgeData(queryData);
      } catch (err2) {
        setErr(err2);
        return;
      }
      setLoaded(true);
    }, []);
    p2(async () => {
      const chainId2 = await window.ethereum.request({ method: "eth_chainId" });
      setChainId(chainId2);
      const handleChainChanged = (chainId3) => {
        setChainId(chainId3);
      };
      window.ethereum.on("chainChanged", handleChainChanged);
    }, []);
    p2(async () => {
      try {
        const accounts = await window.ethereum.request({ method: "eth_requestAccounts" });
        setAccount(accounts[0]);
      } catch (err2) {
        setErr(err2);
      }
      await window.ethereum.request({ method: "eth_accounts" });
      const handleAccountsChanged = (accounts) => {
        setAccount(accounts[0]);
      };
      window.ethereum.on("accountsChanged", handleAccountsChanged);
    }, []);
    p2(async () => {
      if (!provider && !bridgeData)
        return;
      const bridgeContract = new Contract(bridgeAddress, bridgeAbi, provider);
      const bridgeFee2 = await bridgeContract.bridgeFee();
      setBridgeFee(bridgeFee2);
    }, [provider]);
    const bridgeIn = T2(async () => {
      if (!provider && !bridgeData)
        return;
      try {
        setErr(null);
        setBridging(true);
        const { walletAddress, amount, symbol } = bridgeData;
        const signer = await provider.getSigner();
        const userBridgeContract = new Contract(bridgeAddress, bridgeAbi, signer);
        const bridgeContract = new Contract(bridgeAddress, bridgeAbi, provider);
        const tokenAddress = await bridgeContract.registeredSymbol(symbol);
        const tokenContract = new Contract(tokenAddress, tokenAbi, provider);
        const tokenDecimals = await tokenContract.decimals();
        const tokenAllowance = await tokenContract.allowance(account, bridgeAddress);
        const amountToBridge = parseUnits(amount.toString(), tokenDecimals);
        const userTokenContract = new Contract(tokenAddress, tokenAbi, signer);
        if (amountToBridge > tokenAllowance) {
          await userTokenContract.approve(bridgeAddress, amountToBridge);
        }
        const options = { value: bridgeFee };
        const tx = { hash: "0x229ea871231f3574d01f9da2007ff2a0ee4652a4cd418a29379427c4f8a89d43" };
        setBridgeHash(tx.hash);
      } catch (err2) {
        if (err2.info && err2.info.error) {
          setErr(err2.info.error);
        } else {
          setErr(err2);
        }
      }
      setBridging(false);
    }, [provider, bridgeFee, bridgeData]);
    const isBridgeDataValid = T2((bridgeData2) => {
      if (!bridgeData2)
        return false;
      if (!bridgeData2.walletAddress || !bridgeData2.symbol || !bridgeData2.amount)
        return false;
      return true;
    }, []);
    return /* @__PURE__ */ y("div", null, /* @__PURE__ */ y(Box, null, /* @__PURE__ */ y(MetaMaskLogoWrap, null, /* @__PURE__ */ y(MetaMaskLogo, null)), /* @__PURE__ */ y(BoxTitle, null, "G45W"), /* @__PURE__ */ y(BoxDescription, null, "Bridge Ethereum to Dero Stargate with Metamask."), (() => {
      if (err) {
        return /* @__PURE__ */ y("div", null, /* @__PURE__ */ y(ErrDiv, null, err.message), /* @__PURE__ */ y(BoxSubmitButton, { onClick: () => location.reload() }, "RELOAD"));
      }
      if (!loaded)
        return;
      if (!isBridgeDataValid(bridgeData)) {
        return /* @__PURE__ */ y(ErrDiv, null, "Missing important bridge data. Close and request a new instance from the application.");
      }
      if (!account) {
        return /* @__PURE__ */ y(WarningDiv, null, "Waiting for user to connect Metamask...");
      }
      if (bridgeHash) {
        return /* @__PURE__ */ y(SuccessDiv, null, "Bridging successful. Check the app and refresh the token interface, your wrapped tokens should appear in a couple of minutes.", /* @__PURE__ */ y("br", null), /* @__PURE__ */ y("br", null), /* @__PURE__ */ y("a", { href: `https://etherscan.io/tx/${bridgeHash}`, target: "_blank", style: "word-break:break-all;" }, bridgeHash));
      }
      return /* @__PURE__ */ y("div", null, /* @__PURE__ */ y(BoxItem, { title: "Bridge Contract", value: bridgeAddress }), /* @__PURE__ */ y(BoxItem, { title: "Ethereum Address", value: account }), /* @__PURE__ */ y(BoxItem, { title: "Dero Address", value: bridgeData.walletAddress }), /* @__PURE__ */ y(BridgeTransfer, null, bridgeData.amount + " " + bridgeData.symbol), /* @__PURE__ */ y(BoxItem, { title: "Bridge In Fee", value: () => {
        return formatEther((bridgeFee || 0).toString()) + " ETH";
      } }), /* @__PURE__ */ y("div", null, /* @__PURE__ */ y(BoxSubmitButton, { onClick: bridgeIn, disabled: bridging }, bridging && /* @__PURE__ */ y(RotateIcon, { className: "material-icons-round" }, "refresh"), bridging ? "BRIDGING" : "BRIDGE IN")));
    })()));
  };
  D(/* @__PURE__ */ y(App, null), document.getElementById("app"));
})();
/*! Bundled license information:

@noble/hashes/esm/utils.js:
  (*! noble-hashes - MIT License (c) 2022 Paul Miller (paulmillr.com) *)

@noble/secp256k1/lib/esm/index.js:
  (*! noble-secp256k1 - MIT License (c) 2019 Paul Miller (paulmillr.com) *)
*/
//# sourceMappingURL=index.js.map
