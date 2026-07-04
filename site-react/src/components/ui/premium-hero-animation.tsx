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

function useThreeCore() {
  const containerRef = useRef<HTMLDivElement>(null);
  const stateRef = useRef<ThreeSceneState | null>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return undefined;

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, 1, 0.1, 1000);
    camera.position.z = 5;

    const renderer = new THREE.WebGLRenderer({ alpha: true, antialias: true, powerPreference: "high-performance", preserveDrawingBuffer: true });
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
      if (!renderer.getContext().isContextLost()) {
        renderer.render(scene, camera);
      }
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
  const threeRef = useThreeCore();

  return (
    <div className={`premium-hero-animation ${className}`.trim()} aria-hidden="true">
      <div ref={threeRef} className="premium-hero-animation__three" />
    </div>
  );
}
