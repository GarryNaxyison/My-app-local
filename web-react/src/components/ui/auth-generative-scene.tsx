import { Canvas, useFrame, useThree } from "@react-three/fiber";
import { useEffect, useMemo, useRef } from "react";
import * as THREE from "three";
import { cn } from "@/lib/utils";

type AuthGenerativeSceneProps = {
  className?: string;
};

const vertexShader = `
  varying vec2 vUv;
  void main() {
    vUv = uv;
    gl_Position = vec4(position.xy, 0.0, 1.0);
  }
`;

const fragmentShader = `
  precision highp float;
  varying vec2 vUv;
  uniform float u_time;
  uniform vec2 u_resolution;

  float hash(vec2 p) {
    p = fract(p * vec2(123.34, 456.21));
    p += dot(p, p + 45.32);
    return fract(p.x * p.y);
  }

  float noise(vec2 p) {
    vec2 i = floor(p);
    vec2 f = fract(p);
    vec2 u = f * f * (3.0 - 2.0 * f);
    float a = hash(i);
    float b = hash(i + vec2(1.0, 0.0));
    float c = hash(i + vec2(0.0, 1.0));
    float d = hash(i + 1.0);
    return mix(mix(a, b, u.x), mix(c, d, u.x), u.y);
  }

  float fbm(vec2 p) {
    float value = 0.0;
    float amp = 0.55;
    for (int i = 0; i < 5; i++) {
      value += amp * noise(p);
      p *= mat2(1.6, -0.42, 0.38, 1.42);
      amp *= 0.52;
    }
    return value;
  }

  void main() {
    vec2 frag = vUv * u_resolution;
    float mn = min(u_resolution.x, u_resolution.y);
    vec2 uv = (frag - 0.5 * u_resolution) / mn;
    uv.x *= 1.2;

    float field = fbm(uv * 2.8 + vec2(u_time * 0.08, -u_time * 0.05));
    float ring = abs(length(uv) - 0.42);
    float orbit = smoothstep(0.09, 0.0, ring + sin(atan(uv.y, uv.x) * 7.0 + u_time * 0.65) * 0.018);
    float aura = smoothstep(0.92, 0.0, length(uv));
    float lattice = smoothstep(0.03, 0.0, abs(fract((uv.x + uv.y + field * 0.16 + u_time * 0.025) * 18.0) - 0.5) - 0.46);

    vec3 base = vec3(0.012, 0.027, 0.065);
    vec3 cyan = vec3(0.17, 0.86, 1.0);
    vec3 mint = vec3(0.47, 1.0, 0.67);
    vec3 gold = vec3(1.0, 0.77, 0.30);
    vec3 color = base;
    color += cyan * aura * 0.42;
    color += mint * orbit * 0.78;
    color += gold * field * 0.18;
    color += cyan * lattice * 0.16;

    float vignette = smoothstep(1.15, 0.22, length(uv));
    gl_FragColor = vec4(color * vignette, 0.96);
  }
`;

function ScenePlane() {
  const material = useRef<THREE.ShaderMaterial | null>(null);
  const { size } = useThree();
  const uniforms = useMemo(
    () => ({
      u_time: { value: 0 },
      u_resolution: { value: new THREE.Vector2(size.width * 2, size.height * 2) },
    }),
    [size.height, size.width],
  );

  useEffect(() => {
    material.current?.uniforms.u_resolution.value.set(size.width * 2, size.height * 2);
  }, [size.height, size.width]);

  useFrame(({ clock }) => {
    if (material.current) material.current.uniforms.u_time.value = clock.getElapsedTime();
  });

  return (
    <mesh>
      <planeGeometry args={[2, 2]} />
      <shaderMaterial ref={material} vertexShader={vertexShader} fragmentShader={fragmentShader} uniforms={uniforms} transparent depthWrite={false} />
    </mesh>
  );
}

export function AuthGenerativeScene({ className }: AuthGenerativeSceneProps) {
  return (
    <div className={cn("auth-generative-scene-v2", className)} aria-hidden="true">
      <Canvas dpr={[1, 1.75]} gl={{ antialias: false, alpha: true }}>
        <ScenePlane />
      </Canvas>
      <div className="auth-generative-scene-v2__shade" />
    </div>
  );
}

export default AuthGenerativeScene;
