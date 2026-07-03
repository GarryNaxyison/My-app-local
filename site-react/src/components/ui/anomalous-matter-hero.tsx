import { Suspense, useEffect, useRef } from "react";
import * as THREE from "three";

type HeroShaderVariant = "dark" | "light";

type GenerativeArtSceneProps = {
  className?: string;
  color?: string;
  particleColor?: string;
  animate?: boolean;
  variant?: HeroShaderVariant;
  transparentBackdrop?: boolean;
};

const shaderPresets = {
  dark: {
    baseColor: [0.0, 0.0, 0.0],
    cloudColor: [0.22, 0.15, 0.08],
    accentBlue: [0.1, 0.33, 0.95],
    accentCyan: [0.05, 0.85, 0.78],
    lineGain: 1.0,
    cloudGain: 1.0,
    backdropAlpha: 1.0,
    particleCount: 90,
    particleSize: 0.015,
    particleOpacity: 0.66,
    meshDiffuse: 0.7,
    meshFresnel: 0.95,
    meshAmbient: 0.18,
    meshOpacity: 0.95,
  },
  light: {
    baseColor: [0.72, 0.84, 0.96],
    cloudColor: [0.42, 0.58, 0.74],
    accentBlue: [0.02, 0.18, 0.56],
    accentCyan: [0.0, 0.55, 0.72],
    lineGain: 1.9,
    cloudGain: 0.72,
    backdropAlpha: 0.92,
    particleCount: 132,
    particleSize: 0.02,
    particleOpacity: 0.82,
    meshDiffuse: 0.9,
    meshFresnel: 1.45,
    meshAmbient: 0.24,
    meshOpacity: 0.98,
  },
} as const;

export function GenerativeArtScene({
  className = "",
  color = "#6bdcff",
  particleColor = "#ffffff",
  animate = false,
  variant = "dark",
  transparentBackdrop = false,
}: GenerativeArtSceneProps) {
  const mountRef = useRef<HTMLDivElement | null>(null);
  const lightRef = useRef<THREE.PointLight | null>(null);

  useEffect(() => {
    const currentMount = mountRef.current;
    if (!currentMount) return undefined;
    const preset = shaderPresets[variant];
    const meshColor = variant === "light" ? "#0a4fc7" : color;
    const pointColor = variant === "light" ? "#0f5fd1" : particleColor;

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(68, currentMount.clientWidth / currentMount.clientHeight, 0.1, 1000);
    camera.position.z = 3.2;

    const renderer = new THREE.WebGLRenderer({
      antialias: false,
      alpha: true,
      powerPreference: "low-power",
    });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 1.15));
    renderer.setSize(currentMount.clientWidth, currentMount.clientHeight);
    renderer.setClearColor(0x000000, 0);
    currentMount.appendChild(renderer.domElement);

    const backdropGeometry = new THREE.PlaneGeometry(2, 2);
    const backdropMaterial = new THREE.ShaderMaterial({
      uniforms: {
        time: { value: 0 },
        resolution: { value: new THREE.Vector2(currentMount.clientWidth, currentMount.clientHeight) },
        baseColor: { value: new THREE.Vector3(...preset.baseColor) },
        cloudColor: { value: new THREE.Vector3(...preset.cloudColor) },
        accentBlue: { value: new THREE.Vector3(...preset.accentBlue) },
        accentCyan: { value: new THREE.Vector3(...preset.accentCyan) },
        lineGain: { value: preset.lineGain },
        cloudGain: { value: preset.cloudGain },
        backdropAlpha: { value: preset.backdropAlpha },
      },
      vertexShader: `
        varying vec2 vUv;
        void main() {
          vUv = uv;
          gl_Position = vec4(position.xy, 0.0, 1.0);
        }
      `,
      fragmentShader: `
        precision highp float;
        uniform vec2 resolution;
        uniform float time;
        uniform vec3 baseColor;
        uniform vec3 cloudColor;
        uniform vec3 accentBlue;
        uniform vec3 accentCyan;
        uniform float lineGain;
        uniform float cloudGain;
        uniform float backdropAlpha;
        varying vec2 vUv;

        float rnd(vec2 p) {
          p = fract(p * vec2(12.9898, 78.233));
          p += dot(p, p + 34.56);
          return fract(p.x * p.y);
        }
        float noise(in vec2 p) {
          vec2 i = floor(p), f = fract(p), u = f * f * (3.0 - 2.0 * f);
          float a = rnd(i), b = rnd(i + vec2(1.0, 0.0)), c = rnd(i + vec2(0.0, 1.0)), d = rnd(i + 1.0);
          return mix(mix(a, b, u.x), mix(c, d, u.x), u.y);
        }
        float fbm(vec2 p) {
          float t = 0.0, a = 1.0;
          mat2 m = mat2(1.0, -0.5, 0.2, 1.2);
          for (int i = 0; i < 3; i++) {
            t += a * noise(p);
            p *= 2.0 * m;
            a *= 0.5;
          }
          return t;
        }
        float clouds(vec2 p) {
          float d = 1.0, t = 0.0;
          for (float i = 0.0; i < 2.0; i++) {
            float a = d * fbm(i * 10.0 + p.x * 0.2 + 0.2 * (1.0 + i) * p.y + d + i * i + p);
            t = mix(t, d, a);
            d = a;
            p *= 2.0 / (i + 1.0);
          }
          return t;
        }
        void main() {
          vec2 frag = vUv * resolution;
          float mn = min(resolution.x, resolution.y);
          vec2 uv = (frag - 0.5 * resolution) / mn;
          vec2 st = uv * vec2(2.0, 1.0);
          vec3 col = baseColor * 0.018;
          float bg = clouds(vec2(st.x + time * 0.26, -st.y));
          uv *= 1.0 - 0.22 * (sin(time * 0.2) * 0.5 + 0.5);
          for (float i = 1.0; i < 7.0; i++) {
            uv += 0.1 * cos(i * vec2(0.1 + 0.01 * i, 0.8) + i * i + time * 0.38 + 0.1 * uv.x);
            vec2 p = uv;
            float d = max(length(p), 0.015);
            col += lineGain * 0.0015 / d * (cos(sin(i) * vec3(1.0, 2.0, 3.0)) + 1.0);
            float b = noise(i + p + bg * 1.731);
            col += lineGain * 0.0022 * b / length(max(abs(p), vec2(b * abs(p.x) * 0.02, abs(p.y))));
            col = mix(col, baseColor + cloudColor * bg * cloudGain, smoothstep(0.0, 1.25, d));
          }
          col += accentBlue * smoothstep(0.65, 0.0, length(st - vec2(-0.44, 0.06))) * 0.18 * lineGain;
          col += accentCyan * smoothstep(0.72, 0.0, length(st - vec2(0.72, -0.08))) * 0.16 * lineGain;
          gl_FragColor = vec4(col, backdropAlpha);
        }
      `,
      depthTest: false,
      depthWrite: false,
      transparent: true,
    });
    const backdrop = new THREE.Mesh(backdropGeometry, backdropMaterial);
    backdrop.renderOrder = -10;
    if (!transparentBackdrop) {
      scene.add(backdrop);
    }

    const particleCount = preset.particleCount;
    const particlePositions = new Float32Array(particleCount * 3);
    for (let i = 0; i < particleCount; i++) {
      particlePositions[i * 3] = (Math.random() - 0.5) * 7.4;
      particlePositions[i * 3 + 1] = (Math.random() - 0.5) * 4.6;
      particlePositions[i * 3 + 2] = -1.5 - Math.random() * 1.8;
    }
    const particlesGeometry = new THREE.BufferGeometry();
    particlesGeometry.setAttribute("position", new THREE.BufferAttribute(particlePositions, 3));
    const particlesMaterial = new THREE.PointsMaterial({
      color: new THREE.Color(pointColor),
      size: preset.particleSize,
      transparent: true,
      opacity: preset.particleOpacity,
      depthWrite: false,
    });
    const particles = new THREE.Points(particlesGeometry, particlesMaterial);
    scene.add(particles);

    const geometry = new THREE.IcosahedronGeometry(1.25, 18);
    const material = new THREE.ShaderMaterial({
      uniforms: {
        time: { value: 0 },
        pointLightPos: { value: new THREE.Vector3(0, 0, 5) },
        color: { value: new THREE.Color(meshColor) },
        diffuseGain: { value: preset.meshDiffuse },
        fresnelGain: { value: preset.meshFresnel },
        ambientGain: { value: preset.meshAmbient },
        meshOpacity: { value: preset.meshOpacity },
      },
      vertexShader: `
        uniform float time;
        varying vec3 vNormal;
        varying vec3 vPosition;

        vec3 mod289(vec3 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
        vec4 mod289(vec4 x) { return x - floor(x * (1.0 / 289.0)) * 289.0; }
        vec4 permute(vec4 x) { return mod289(((x * 34.0) + 1.0) * x); }
        vec4 taylorInvSqrt(vec4 r) { return 1.79284291400159 - 0.85373472095314 * r; }

        float snoise(vec3 v) {
          const vec2 C = vec2(1.0 / 6.0, 1.0 / 3.0);
          const vec4 D = vec4(0.0, 0.5, 1.0, 2.0);
          vec3 i = floor(v + dot(v, C.yyy));
          vec3 x0 = v - i + dot(i, C.xxx);
          vec3 g = step(x0.yzx, x0.xyz);
          vec3 l = 1.0 - g;
          vec3 i1 = min(g.xyz, l.zxy);
          vec3 i2 = max(g.xyz, l.zxy);
          vec3 x1 = x0 - i1 + C.xxx;
          vec3 x2 = x0 - i2 + C.yyy;
          vec3 x3 = x0 - D.yyy;
          i = mod289(i);
          vec4 p = permute(permute(permute(
                    i.z + vec4(0.0, i1.z, i2.z, 1.0))
                  + i.y + vec4(0.0, i1.y, i2.y, 1.0))
                  + i.x + vec4(0.0, i1.x, i2.x, 1.0));
          float n_ = 0.142857142857;
          vec3 ns = n_ * D.wyz - D.xzx;
          vec4 j = p - 49.0 * floor(p * ns.z * ns.z);
          vec4 x_ = floor(j * ns.z);
          vec4 y_ = floor(j - 7.0 * x_);
          vec4 x = x_ * ns.x + ns.yyyy;
          vec4 y = y_ * ns.x + ns.yyyy;
          vec4 h = 1.0 - abs(x) - abs(y);
          vec4 b0 = vec4(x.xy, y.xy);
          vec4 b1 = vec4(x.zw, y.zw);
          vec4 s0 = floor(b0) * 2.0 + 1.0;
          vec4 s1 = floor(b1) * 2.0 + 1.0;
          vec4 sh = -step(h, vec4(0.0));
          vec4 a0 = b0.xzyw + s0.xzyw * sh.xxyy;
          vec4 a1 = b1.xzyw + s1.xzyw * sh.zzww;
          vec3 p0 = vec3(a0.xy, h.x);
          vec3 p1 = vec3(a0.zw, h.y);
          vec3 p2 = vec3(a1.xy, h.z);
          vec3 p3 = vec3(a1.zw, h.w);
          vec4 norm = taylorInvSqrt(vec4(dot(p0, p0), dot(p1, p1), dot(p2, p2), dot(p3, p3)));
          p0 *= norm.x;
          p1 *= norm.y;
          p2 *= norm.z;
          p3 *= norm.w;
          vec4 m = max(0.6 - vec4(dot(x0, x0), dot(x1, x1), dot(x2, x2), dot(x3, x3)), 0.0);
          m = m * m;
          return 42.0 * dot(m * m, vec4(dot(p0, x0), dot(p1, x1), dot(p2, x2), dot(p3, x3)));
        }

        void main() {
          vNormal = normalize(normalMatrix * normal);
          vPosition = position;
          float displacement = snoise(position * 2.0 + time * 0.5) * 0.22;
          vec3 newPosition = position + normal * displacement;
          gl_Position = projectionMatrix * modelViewMatrix * vec4(newPosition, 1.0);
        }
      `,
      fragmentShader: `
        uniform vec3 color;
        uniform vec3 pointLightPos;
        uniform float diffuseGain;
        uniform float fresnelGain;
        uniform float ambientGain;
        uniform float meshOpacity;
        varying vec3 vNormal;
        varying vec3 vPosition;

        void main() {
          vec3 normal = normalize(vNormal);
          vec3 lightDir = normalize(pointLightPos - vPosition);
          float diffuse = max(dot(normal, lightDir), 0.0);
          float fresnel = pow(1.0 - max(dot(normal, vec3(0.0, 0.0, 1.0)), 0.0), 2.0);
          vec3 finalColor = color * (diffuse * diffuseGain + fresnel * fresnelGain + ambientGain);
          gl_FragColor = vec4(finalColor, meshOpacity);
        }
      `,
      wireframe: true,
      transparent: true,
    });

    const mesh = new THREE.Mesh(geometry, material);
    mesh.rotation.z = -0.22;
    scene.add(mesh);

    const pointLight = new THREE.PointLight(0xffffff, 1, 100);
    pointLight.position.set(0, 0, 5);
    lightRef.current = pointLight;
    scene.add(pointLight);

    let frameId = 0;
    let lastFrame = 0;
    let isVisible = true;
    let isDocumentVisible = !document.hidden;
    const reducedMotionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
    let isReducedMotion = reducedMotionQuery.matches;
    let isLoopRunning = false;
    const renderScene = (time = 0) => {
      const seconds = time * 0.001;
      backdropMaterial.uniforms.time.value = seconds;
      material.uniforms.time.value = time * 0.00036;
      mesh.rotation.y = -0.32 + seconds * 0.12;
      mesh.rotation.x = seconds * 0.04;
      particles.rotation.z = seconds * 0.012;
      particles.rotation.y = seconds * 0.006;
      renderer.render(scene, camera);
    };
    const shouldAnimate = () => animate && isVisible && isDocumentVisible && !isReducedMotion;
    const runAnimation = (time: number) => {
      if (!shouldAnimate()) {
        isLoopRunning = false;
        frameId = 0;
        return;
      }
      if (time - lastFrame >= 50) {
        renderScene(time);
        lastFrame = time;
      }
      frameId = requestAnimationFrame(runAnimation);
    };
    const startLoop = () => {
      if (!shouldAnimate() || isLoopRunning) return;
      isLoopRunning = true;
      frameId = requestAnimationFrame(runAnimation);
    };
    const stopLoop = () => {
      if (frameId) cancelAnimationFrame(frameId);
      frameId = 0;
      isLoopRunning = false;
    };
    const syncLoop = () => {
      if (shouldAnimate()) {
        startLoop();
        return;
      }
      stopLoop();
      renderScene(performance.now());
    };

    syncLoop();

    const handleResize = () => {
      camera.aspect = currentMount.clientWidth / currentMount.clientHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(currentMount.clientWidth, currentMount.clientHeight);
      backdropMaterial.uniforms.resolution.value.set(currentMount.clientWidth, currentMount.clientHeight);
      renderScene(performance.now());
    };

    const handleMouseMove = (event: MouseEvent) => {
      const x = (event.clientX / window.innerWidth) * 2 - 1;
      const y = -(event.clientY / window.innerHeight) * 2 + 1;
      const vector = new THREE.Vector3(x, y, 0.5).unproject(camera);
      const direction = vector.sub(camera.position).normalize();
      const distance = -camera.position.z / direction.z;
      const position = camera.position.clone().add(direction.multiplyScalar(distance));
      lightRef.current?.position.copy(position);
      material.uniforms.pointLightPos.value = position;
      renderScene(performance.now());
    };
    const handleVisibilityChange = () => {
      isDocumentVisible = !document.hidden;
      syncLoop();
    };
    const handleReducedMotionChange = () => {
      isReducedMotion = reducedMotionQuery.matches;
      syncLoop();
    };
    const observer = new IntersectionObserver(
      ([entry]) => {
        isVisible = Boolean(entry?.isIntersecting);
        syncLoop();
      },
      { threshold: 0.02 },
    );
    observer.observe(currentMount);

    window.addEventListener("resize", handleResize);
    window.addEventListener("mousemove", handleMouseMove);
    document.addEventListener("visibilitychange", handleVisibilityChange);
    reducedMotionQuery.addEventListener("change", handleReducedMotionChange);

    return () => {
      stopLoop();
      observer.disconnect();
      window.removeEventListener("resize", handleResize);
      window.removeEventListener("mousemove", handleMouseMove);
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      reducedMotionQuery.removeEventListener("change", handleReducedMotionChange);
      geometry.dispose();
      material.dispose();
      backdropGeometry.dispose();
      backdropMaterial.dispose();
      particlesGeometry.dispose();
      particlesMaterial.dispose();
      renderer.dispose();
      if (renderer.domElement.parentNode === currentMount) {
        currentMount.removeChild(renderer.domElement);
      }
    };
  }, [animate, color, particleColor, transparentBackdrop, variant]);

  return <div ref={mountRef} className={`generative-art-scene ${className}`} data-hero-preset={variant} aria-hidden="true" />;
}

type AnomalousMatterHeroProps = {
  title?: string;
  subtitle?: string;
  description?: string;
};

export function AnomalousMatterHero({
  title = "NERIVA",
  subtitle = "Язык, который тренируется каждый день",
  description = "Уроки, роли, произношение, заметки и Telegram в одном учебном профиле.",
}: AnomalousMatterHeroProps) {
  return (
    <section role="banner" className="anomalous-matter-hero">
      <Suspense fallback={<div className="anomalous-matter-hero__fallback" />}>
        <GenerativeArtScene animate />
      </Suspense>
      <div className="anomalous-matter-hero__shade" />
      <div className="anomalous-matter-hero__content">
        <h1>{title}</h1>
        <p>{subtitle}</p>
        <span>{description}</span>
      </div>
    </section>
  );
}
