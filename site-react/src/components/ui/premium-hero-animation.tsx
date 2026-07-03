import { useEffect, useRef } from "react";
import * as THREE from "three";

type ThreeSceneState = {
  renderer: THREE.WebGLRenderer;
  scene: THREE.Scene;
  camera: THREE.PerspectiveCamera;
  coreGroup: THREE.Group;
  core: THREE.Mesh;
  innerSphere: THREE.Mesh;
  particles: THREE.Points;
  frame: number;
  resizeObserver?: ResizeObserver;
};

function useFluidShader() {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    const canvas = canvasRef.current;
    const gl = (canvas?.getContext("webgl") || canvas?.getContext("experimental-webgl")) as WebGLRenderingContext | null;
    if (!canvas || !gl) return undefined;

    const vertexShader = `
      attribute vec2 a_position;
      varying vec2 v_texCoord;
      void main() {
        v_texCoord = a_position * 0.5 + 0.5;
        gl_Position = vec4(a_position, 0.0, 1.0);
      }
    `;

    const fragmentShader = `
      precision highp float;
      uniform float u_time;
      uniform vec2 u_resolution;
      uniform vec2 u_mouse;
      varying vec2 v_texCoord;

      void main() {
        vec2 uv = v_texCoord;
        vec2 p = (v_texCoord - 0.5) * 2.0;
        p.x *= u_resolution.x / u_resolution.y;

        float t = u_time * 0.2;
        float wave = sin(p.x * 2.0 + t) * 0.5 + 0.5;
        wave += sin(p.y * 3.0 - t * 1.5) * 0.3;

        vec3 color1 = vec3(0.043, 0.078, 0.133);
        vec3 color2 = vec3(0.0, 0.196, 0.627);
        vec3 color3 = vec3(0.239, 0.482, 1.0);

        float noise = sin(p.x * 4.0 + sin(t + p.y * 2.0)) * 0.5 + 0.5;
        vec3 finalColor = mix(color1, color2, wave * 0.6);
        finalColor = mix(finalColor, color3, pow(noise, 8.0) * 0.15);

        float dist = length(uv - u_mouse / u_resolution);
        finalColor += color3 * (1.0 - smoothstep(0.0, 0.4, dist)) * 0.05;

        gl_FragColor = vec4(finalColor, 1.0);
      }
    `;

    const compile = (type: number, source: string) => {
      const shader = gl.createShader(type);
      if (!shader) return null;
      gl.shaderSource(shader, source);
      gl.compileShader(shader);
      if (!gl.getShaderParameter(shader, gl.COMPILE_STATUS)) {
        gl.deleteShader(shader);
        return null;
      }
      return shader;
    };

    const vs = compile(gl.VERTEX_SHADER, vertexShader);
    const fs = compile(gl.FRAGMENT_SHADER, fragmentShader);
    const program = gl.createProgram();
    const buffer = gl.createBuffer();
    if (!vs || !fs || !program || !buffer) return undefined;

    gl.attachShader(program, vs);
    gl.attachShader(program, fs);
    gl.linkProgram(program);
    gl.useProgram(program);
    gl.bindBuffer(gl.ARRAY_BUFFER, buffer);
    gl.bufferData(gl.ARRAY_BUFFER, new Float32Array([-1, -1, 1, -1, -1, 1, 1, 1]), gl.STATIC_DRAW);

    const position = gl.getAttribLocation(program, "a_position");
    gl.enableVertexAttribArray(position);
    gl.vertexAttribPointer(position, 2, gl.FLOAT, false, 0, 0);

    const uTime = gl.getUniformLocation(program, "u_time");
    const uResolution = gl.getUniformLocation(program, "u_resolution");
    const uMouse = gl.getUniformLocation(program, "u_mouse");
    const mouse = { x: 0, y: 0 };
    let frame = 0;

    const syncSize = () => {
      const dpr = Math.max(1, Math.min(window.devicePixelRatio || 1, 2));
      const rect = canvas.getBoundingClientRect();
      const width = Math.max(1, Math.floor(rect.width * dpr));
      const height = Math.max(1, Math.floor(rect.height * dpr));
      if (canvas.width !== width || canvas.height !== height) {
        canvas.width = width;
        canvas.height = height;
      }
      if (!mouse.x && !mouse.y) {
        mouse.x = width * 0.56;
        mouse.y = height * 0.58;
      }
    };

    const onMouseMove = (event: MouseEvent) => {
      const rect = canvas.getBoundingClientRect();
      if (!rect.width || !rect.height) return;
      const nx = (event.clientX - rect.left) / rect.width;
      const ny = 1 - (event.clientY - rect.top) / rect.height;
      mouse.x = nx * canvas.width;
      mouse.y = ny * canvas.height;
    };

    const resizeObserver = typeof ResizeObserver !== "undefined" ? new ResizeObserver(syncSize) : undefined;
    resizeObserver?.observe(canvas);
    syncSize();
    window.addEventListener("mousemove", onMouseMove);

    const render = (time: number) => {
      syncSize();
      gl.viewport(0, 0, canvas.width, canvas.height);
      gl.useProgram(program);
      if (uTime) gl.uniform1f(uTime, time * 0.001);
      if (uResolution) gl.uniform2f(uResolution, canvas.width, canvas.height);
      if (uMouse) gl.uniform2f(uMouse, mouse.x, mouse.y);
      gl.drawArrays(gl.TRIANGLE_STRIP, 0, 4);
      frame = requestAnimationFrame(render);
    };

    frame = requestAnimationFrame(render);

    return () => {
      resizeObserver?.disconnect();
      window.removeEventListener("mousemove", onMouseMove);
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

function useThreeCore() {
  const containerRef = useRef<HTMLDivElement>(null);
  const stateRef = useRef<ThreeSceneState | null>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return undefined;

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, 1, 0.1, 1000);
    camera.position.z = 5;

    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true });
    renderer.setClearColor(0x000000, 0);
    renderer.setPixelRatio(Math.min(window.devicePixelRatio || 1, 2));
    container.appendChild(renderer.domElement);

    const coreGroup = new THREE.Group();
    scene.add(coreGroup);

    const coreGeometry = new THREE.IcosahedronGeometry(1.5, 4);
    const coreMaterial = new THREE.MeshPhongMaterial({
      color: 0x3d7bff,
      emissive: 0x0032a0,
      wireframe: true,
      transparent: true,
      opacity: 0.42,
    });
    const core = new THREE.Mesh(coreGeometry, coreMaterial);
    coreGroup.add(core);

    const innerGeometry = new THREE.SphereGeometry(1, 32, 32);
    const innerMaterial = new THREE.MeshStandardMaterial({
      color: 0x0032a0,
      emissive: 0x3d7bff,
      emissiveIntensity: 0.5,
      roughness: 0.2,
      metalness: 0.8,
      transparent: true,
      opacity: 0.82,
    });
    const innerSphere = new THREE.Mesh(innerGeometry, innerMaterial);
    coreGroup.add(innerSphere);

    const particlesGeometry = new THREE.BufferGeometry();
    const particlesCount = 1000;
    const positions = new Float32Array(particlesCount * 3);
    for (let index = 0; index < positions.length; index += 1) {
      positions[index] = (Math.random() - 0.5) * 10;
    }
    particlesGeometry.setAttribute("position", new THREE.BufferAttribute(positions, 3));
    const particlesMaterial = new THREE.PointsMaterial({
      size: 0.02,
      color: 0xffffff,
      transparent: true,
      opacity: 0.8,
    });
    const particles = new THREE.Points(particlesGeometry, particlesMaterial);
    scene.add(particles);

    scene.add(new THREE.AmbientLight(0xffffff, 0.5));
    const pointLight = new THREE.PointLight(0x3d7bff, 2);
    pointLight.position.set(5, 5, 5);
    scene.add(pointLight);

    const resize = () => {
      const width = Math.max(1, container.clientWidth);
      const height = Math.max(1, container.clientHeight);
      camera.aspect = width / height;
      camera.updateProjectionMatrix();
      renderer.setSize(width, height, false);
    };

    const resizeObserver = typeof ResizeObserver !== "undefined" ? new ResizeObserver(resize) : undefined;
    resizeObserver?.observe(container);
    resize();

    const animate = () => {
      const time = Date.now() * 0.002;
      coreGroup.rotation.y += 0.005;
      coreGroup.rotation.z += 0.002;
      particles.rotation.y += 0.001;
      core.scale.setScalar(1 + Math.sin(time) * 0.05);
      renderer.render(scene, camera);
      if (stateRef.current) stateRef.current.frame = requestAnimationFrame(animate);
    };

    stateRef.current = { renderer, scene, camera, coreGroup, core, innerSphere, particles, frame: 0, resizeObserver };
    animate();

    return () => {
      const state = stateRef.current;
      if (!state) return;
      cancelAnimationFrame(state.frame);
      state.resizeObserver?.disconnect();
      if (state.renderer.domElement.parentElement === container) {
        container.removeChild(state.renderer.domElement);
      }
      state.renderer.dispose();
      (core.geometry as THREE.BufferGeometry).dispose();
      (core.material as THREE.Material).dispose();
      (innerSphere.geometry as THREE.BufferGeometry).dispose();
      (innerSphere.material as THREE.Material).dispose();
      (particles.geometry as THREE.BufferGeometry).dispose();
      (particles.material as THREE.Material).dispose();
      stateRef.current = null;
    };
  }, []);

  return containerRef;
}

export function PremiumHeroAnimation({ className = "" }: { className?: string }) {
  const shaderRef = useFluidShader();
  const threeRef = useThreeCore();

  return (
    <div className={`premium-hero-animation ${className}`.trim()} aria-hidden="true">
      <canvas ref={shaderRef} className="premium-hero-animation__shader" />
      <div ref={threeRef} className="premium-hero-animation__three" />
    </div>
  );
}
