import { useEffect, useRef } from "react";

type HeroProps = {
  trustBadge?: {
    text: string;
    icons?: string[];
  };
  headline: {
    line1: string;
    line2: string;
  };
  subtitle: string;
  buttons?: {
    primary?: {
      text: string;
      href?: string;
      onClick?: () => void;
    };
    secondary?: {
      text: string;
      href?: string;
      onClick?: () => void;
    };
  };
  className?: string;
};

const defaultShaderSource = `#version 300 es
precision highp float;
out vec4 O;
uniform vec2 resolution;
uniform float time;
#define FC gl_FragCoord.xy
#define T time
#define R resolution
#define MN min(R.x,R.y)
float rnd(vec2 p) {
  p = fract(p * vec2(12.9898, 78.233));
  p += dot(p, p + 34.56);
  return fract(p.x * p.y);
}
float noise(in vec2 p) {
  vec2 i = floor(p), f = fract(p), u = f * f * (3. - 2. * f);
  float a = rnd(i), b = rnd(i + vec2(1, 0)), c = rnd(i + vec2(0, 1)), d = rnd(i + 1.);
  return mix(mix(a, b, u.x), mix(c, d, u.x), u.y);
}
float fbm(vec2 p) {
  float t = .0, a = 1.;
  mat2 m = mat2(1., -.5, .2, 1.2);
  for (int i = 0; i < 5; i++) {
    t += a * noise(p);
    p *= 2. * m;
    a *= .5;
  }
  return t;
}
float clouds(vec2 p) {
  float d = 1., t = .0;
  for (float i = .0; i < 3.; i++) {
    float a = d * fbm(i * 10. + p.x * .2 + .2 * (1. + i) * p.y + d + i * i + p);
    t = mix(t, d, a);
    d = a;
    p *= 2. / (i + 1.);
  }
  return t;
}
void main(void) {
  vec2 uv = (FC - .5 * R) / MN, st = uv * vec2(2, 1);
  vec3 col = vec3(0);
  float bg = clouds(vec2(st.x + T * .5, -st.y));
  uv *= 1. - .3 * (sin(T * .2) * .5 + .5);
  for (float i = 1.; i < 12.; i++) {
    uv += .1 * cos(i * vec2(.1 + .01 * i, .8) + i * i + T * .5 + .1 * uv.x);
    vec2 p = uv;
    float d = length(p);
    col += .00125 / d * (cos(sin(i) * vec3(1, 2, 3)) + 1.);
    float b = noise(i + p + bg * 1.731);
    col += .002 * b / length(max(p, vec2(b * p.x * .02, p.y)));
    col = mix(col, vec3(bg * .25, bg * .137, bg * .05), d);
  }
  O = vec4(col, 1);
}`;

function useShaderBackground() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    const gl = canvas?.getContext("webgl2");
    if (!canvas || !gl) return undefined;

    const vertexSrc = `#version 300 es
      precision highp float;
      in vec4 position;
      void main(){ gl_Position = position; }
    `;
    const vertices = new Float32Array([-1, 1, -1, -1, 1, 1, 1, -1]);

    const compile = (type: number, source: string) => {
      const shader = gl.createShader(type);
      if (!shader) throw new Error("Could not create shader");
      gl.shaderSource(shader, source);
      gl.compileShader(shader);
      if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
        throw new Error(gl.getShaderInfoLog(shader) || "Shader compile failed");
      }
      return shader;
    };

    let frame = 0;
    const program = gl.createProgram();
    const buffer = gl.createBuffer();
    const vs = compile(gl.VERTEX_SHADER, vertexSrc);
    const fs = compile(gl.FRAGMENT_SHADER, defaultShaderSource);
    if (!program || !buffer) return undefined;

    gl.attachShader(program, vs);
    gl.attachShader(program, fs);
    gl.linkProgram(program);
    gl.useProgram(program);

    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW);
    const position = gl.getAttribLocation(program, "position");
    gl.enableVertexAttribArray(position);
    gl.vertexAttribPointer(position, 2, gl.FLOAT, false, 0, 0);

    const resolution = gl.getUniformLocation(program, "resolution");
    const time = gl.getUniformLocation(program, "time");

    const resize = () => {
      const dpr = Math.max(1, Math.min(window.devicePixelRatio, 2));
      const rect = canvas.getBoundingClientRect();
      canvas.width = Math.max(1, Math.floor(rect.width * dpr));
      canvas.height = Math.max(1, Math.floor(rect.height * dpr));
      gl.viewport(0, 0, canvas.width, canvas.height);
    };

    const loop = (now: number) => {
      gl.clearColor(0, 0, 0, 1);
      gl.clear(gl.COLOR_BUFFER_BIT);
      gl.useProgram(program);
      gl.uniform2f(resolution, canvas.width, canvas.height);
      gl.uniform1f(time, now * 0.001);
      gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
      frame = requestAnimationFrame(loop);
    };

    resize();
    window.addEventListener("resize", resize);
    frame = requestAnimationFrame(loop);

    return () => {
      window.removeEventListener("resize", resize);
      cancelAnimationFrame(frame);
      gl.deleteBuffer(buffer);
      gl.detachShader(program, vs);
      gl.detachShader(program, fs);
      gl.deleteShader(vs);
      gl.deleteShader(fs);
      gl.deleteProgram(program);
    };
  }, []);

  return canvasRef;
}

export default function AnimatedShaderHero({ trustBadge, headline, subtitle, buttons, className = "" }: HeroProps) {
  const canvasRef = useShaderBackground();
  const buttonMarkup = (kind: "primary" | "secondary", button?: { text: string; href?: string; onClick?: () => void }) => {
    if (!button) return null;
    const className = kind === "primary" ? "hero-action hero-action--primary" : "hero-action hero-action--secondary";
    if (button.href) {
      return (
        <a className={className} href={button.href} onClick={button.onClick}>
          {button.text}
        </a>
      );
    }
    return (
      <button className={className} type="button" onClick={button.onClick}>
        {button.text}
      </button>
    );
  };

  return (
    <section className={`shader-hero ${className}`}>
      <canvas ref={canvasRef} className="shader-hero__canvas" />
      <div className="shader-hero__veil" />
      <div className="shader-hero__content">
        {trustBadge ? (
          <div className="trust-badge">
            <span>{trustBadge.icons?.join(" ") || "*"}</span>
            {trustBadge.text}
          </div>
        ) : null}
        <h1>
          <span>{headline.line1}</span>
          <span>{headline.line2}</span>
        </h1>
        <p>{subtitle}</p>
        <div className="hero-actions">
          {buttonMarkup("primary", buttons?.primary)}
          {buttonMarkup("secondary", buttons?.secondary)}
        </div>
      </div>
    </section>
  );
}
