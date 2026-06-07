import { Canvas, useFrame, useThree } from "@react-three/fiber";
import { useEffect, useMemo, useRef } from "react";
import * as THREE from "three";
import { cn } from "../lib/format";

type CanvasRevealEffectProps = {
  className?: string;
  colors?: [number, number, number][];
  dotSize?: number;
  animationSpeed?: number;
  reverse?: boolean;
  showGradient?: boolean;
};

function rgb(color: [number, number, number]) {
  return new THREE.Color(color[0] / 255, color[1] / 255, color[2] / 255);
}

const vertexShader = `
  varying vec2 vUv;

  void main() {
    vUv = uv;
    gl_Position = vec4(position.xy, 0.0, 1.0);
  }
`;

const fragmentShader = `
  precision mediump float;

  varying vec2 vUv;
  uniform float u_time;
  uniform vec2 u_resolution;
  uniform vec3 u_color_a;
  uniform vec3 u_color_b;
  uniform float u_dot_size;
  uniform float u_speed;
  uniform float u_reverse;

  float random(vec2 st) {
    return fract(sin(dot(st.xy, vec2(12.9898,78.233))) * 43758.5453123);
  }

  void main() {
    vec2 pixel = vUv * u_resolution;
    float cell = max(u_dot_size * 5.0, 18.0);
    vec2 grid = floor(pixel / cell);
    vec2 local = fract(pixel / cell) - 0.5;
    float dot = 1.0 - smoothstep(0.10, 0.18, length(local));

    vec2 centered = vUv * 2.0 - 1.0;
    centered.x *= u_resolution.x / max(u_resolution.y, 1.0);
    float distance_from_center = length(centered);
    float wave = u_time * u_speed * 0.18;
    float reveal = smoothstep(wave - 0.42, wave, distance_from_center);
    if (u_reverse > 0.5) {
      reveal = 1.0 - reveal;
    }

    float noise = random(grid + floor(u_time * 2.0));
    float shimmer = 0.55 + 0.45 * sin(u_time * 1.7 + random(grid) * 6.28318);
    float field = dot * reveal * mix(0.22, 1.0, noise) * shimmer;

    vec3 color = mix(u_color_a, u_color_b, vUv.x + random(grid) * 0.18);
    gl_FragColor = vec4(color * field, field);
  }
`;

function ShaderPlane({
  colors,
  dotSize,
  animationSpeed,
  reverse,
}: Required<Pick<CanvasRevealEffectProps, "colors" | "dotSize" | "animationSpeed" | "reverse">>) {
  const material = useRef<THREE.ShaderMaterial>(null);
  const { size } = useThree();

  const uniforms = useMemo(
    () => ({
      u_time: { value: 0 },
      u_resolution: { value: new THREE.Vector2(size.width * 2, size.height * 2) },
      u_color_a: { value: rgb(colors[0]) },
      u_color_b: { value: rgb(colors[1] ?? colors[0]) },
      u_dot_size: { value: dotSize },
      u_speed: { value: animationSpeed },
      u_reverse: { value: reverse ? 1 : 0 },
    }),
    [animationSpeed, colors, dotSize, reverse, size.height, size.width],
  );

  useEffect(() => {
    if (!material.current) return;
    material.current.uniforms.u_resolution.value.set(size.width * 2, size.height * 2);
    material.current.uniforms.u_color_a.value = rgb(colors[0]);
    material.current.uniforms.u_color_b.value = rgb(colors[1] ?? colors[0]);
    material.current.uniforms.u_reverse.value = reverse ? 1 : 0;
  }, [colors, reverse, size.height, size.width]);

  useFrame(({ clock }) => {
    if (!material.current) return;
    material.current.uniforms.u_time.value = clock.getElapsedTime();
  });

  return (
    <mesh>
      <planeGeometry args={[2, 2]} />
      <shaderMaterial
        ref={material}
        vertexShader={vertexShader}
        fragmentShader={fragmentShader}
        uniforms={uniforms}
        transparent
        depthWrite={false}
        blending={THREE.AdditiveBlending}
      />
    </mesh>
  );
}

export function CanvasRevealEffect({
  className,
  colors = [
    [43, 212, 255],
    [255, 214, 102],
  ],
  dotSize = 4,
  animationSpeed = 3,
  reverse = false,
  showGradient = true,
}: CanvasRevealEffectProps) {
  return (
    <div className={cn("canvas-reveal", className)}>
      <Canvas dpr={[1, 1.75]} gl={{ antialias: false, alpha: true }}>
        <ShaderPlane colors={colors} dotSize={dotSize} animationSpeed={animationSpeed} reverse={reverse} />
      </Canvas>
      {showGradient ? <div className="canvas-reveal__shade" /> : null}
    </div>
  );
}
