"use client";

import { useMemo, useRef } from "react";
import { useFrame } from "@react-three/fiber";
import { DotOrbit, MeshGradient } from "@paper-design/shaders-react";
import * as THREE from "three";
import { cn } from "@/lib/utils";

const vertexShader = `
  uniform float time;
  uniform float intensity;
  varying vec2 vUv;
  varying vec3 vPosition;

  void main() {
    vUv = uv;
    vPosition = position;

    vec3 pos = position;
    pos.y += sin(pos.x * 10.0 + time) * 0.1 * intensity;
    pos.x += cos(pos.y * 8.0 + time * 1.5) * 0.05 * intensity;

    gl_Position = projectionMatrix * modelViewMatrix * vec4(pos, 1.0);
  }
`;

const fragmentShader = `
  uniform float time;
  uniform float intensity;
  uniform vec3 color1;
  uniform vec3 color2;
  varying vec2 vUv;
  varying vec3 vPosition;

  void main() {
    vec2 uv = vUv;
    float noise = sin(uv.x * 20.0 + time) * cos(uv.y * 15.0 + time * 0.8);
    noise += sin(uv.x * 35.0 - time * 2.0) * cos(uv.y * 25.0 + time * 1.2) * 0.5;

    vec3 color = mix(color1, color2, noise * 0.5 + 0.5);
    color = mix(color, vec3(1.0), pow(abs(noise), 2.0) * intensity);

    float glow = 1.0 - length(uv - 0.5) * 2.0;
    glow = pow(glow, 2.0);

    gl_FragColor = vec4(color * glow, glow * 0.8);
  }
`;

export function ShaderPlane({
  position,
  color1 = "#ff5722",
  color2 = "#ffffff",
}: {
  position: [number, number, number];
  color1?: string;
  color2?: string;
}) {
  const mesh = useRef<THREE.Mesh>(null);

  const uniforms = useMemo(
    () => ({
      time: { value: 0 },
      intensity: { value: 1.0 },
      color1: { value: new THREE.Color(color1) },
      color2: { value: new THREE.Color(color2) },
    }),
    [color1, color2],
  );

  useFrame((state) => {
    if (!mesh.current) return;
    uniforms.time.value = state.clock.elapsedTime;
    uniforms.intensity.value = 1.0 + Math.sin(state.clock.elapsedTime * 2) * 0.3;
  });

  return (
    <mesh ref={mesh} position={position}>
      <planeGeometry args={[2, 2, 32, 32]} />
      <shaderMaterial uniforms={uniforms} vertexShader={vertexShader} fragmentShader={fragmentShader} transparent side={THREE.DoubleSide} />
    </mesh>
  );
}

export function EnergyRing({
  radius = 1,
  position = [0, 0, 0],
}: {
  radius?: number;
  position?: [number, number, number];
}) {
  const mesh = useRef<THREE.Mesh>(null);

  useFrame((state) => {
    if (!mesh.current) return;
    mesh.current.rotation.z = state.clock.elapsedTime;
    const material = mesh.current.material as THREE.MeshBasicMaterial;
    material.opacity = 0.5 + Math.sin(state.clock.elapsedTime * 3) * 0.3;
  });

  return (
    <mesh ref={mesh} position={position}>
      <ringGeometry args={[radius * 0.8, radius, 32]} />
      <meshBasicMaterial color="#ff5722" transparent opacity={0.6} side={THREE.DoubleSide} />
    </mesh>
  );
}

export function PaperShaderBackground({ theme, className }: { theme: "light" | "dark"; className?: string }) {
  const dark = theme === "dark";
  return (
    <div className={cn("paper-shader-bg", dark ? "paper-shader-bg--dark" : "paper-shader-bg--light", className)}>
      <MeshGradient
        className="paper-shader-bg__mesh"
        colors={dark ? ["#020511", "#0b1932", "#114b58", "#7aa2ff"] : ["#f8fbff", "#d7ecff", "#b9f3e9", "#ffe8b8"]}
        speed={dark ? 0.62 : 0.5}
        distortion={dark ? 0.86 : 0.64}
        swirl={dark ? 0.24 : 0.18}
      />
      <DotOrbit
        className="paper-shader-bg__orbit"
        colors={dark ? ["#7aa2ff", "#45d8bb", "#f4c96d"] : ["#2d5bff", "#45d8bb", "#f3b84b"]}
        colorBack={dark ? "#020511" : "#edf7ff"}
        speed={dark ? 0.72 : 0.58}
        size={dark ? 0.28 : 0.22}
        sizeRange={dark ? 0.24 : 0.16}
        scale={dark ? 0.72 : 0.64}
      />
      <div className="paper-shader-bg__lights" />
    </div>
  );
}
