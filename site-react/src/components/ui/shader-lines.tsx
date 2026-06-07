"use client";

import { useEffect, useRef } from "react";
import * as THREE from "three";

export function ShaderLines({ className = "" }: { className?: string }) {
  const containerRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<{
    camera: THREE.Camera;
    scene: THREE.Scene;
    renderer: THREE.WebGLRenderer;
    uniforms: { time: { value: number }; resolution: { value: THREE.Vector2 } };
    animationId: number;
    geometry: THREE.PlaneGeometry;
    material: THREE.ShaderMaterial;
  } | null>(null);

  useEffect(() => {
    const container = containerRef.current;
    if (!container) return undefined;

    const vertexShader = `
      void main() {
        gl_Position = vec4(position, 1.0);
      }
    `;

    const fragmentShader = `
      precision highp float;
      uniform vec2 resolution;
      uniform float time;

      float random(float x) {
        return fract(sin(x) * 1e4);
      }

      void main(void) {
        vec2 uv = (gl_FragCoord.xy * 2.0 - resolution.xy) / min(resolution.x, resolution.y);
        vec2 mosaic = vec2(4.0, 2.0);
        vec2 screen = vec2(256.0, 256.0);
        uv.x = floor(uv.x * screen.x / mosaic.x) / (screen.x / mosaic.x);
        uv.y = floor(uv.y * screen.y / mosaic.y) / (screen.y / mosaic.y);

        float t = time * 0.06 + random(uv.x) * 0.4;
        float lineWidth = 0.0008;
        vec3 color = vec3(0.0);

        for (int j = 0; j < 3; j++) {
          for (int i = 0; i < 5; i++) {
            color[j] += lineWidth * float(i * i) / abs(fract(t - 0.01 * float(j) + float(i) * 0.01) - length(uv));
          }
        }

        gl_FragColor = vec4(color[2], color[1], color[0], 1.0);
      }
    `;

    const camera = new THREE.Camera();
    camera.position.z = 1;
    const scene = new THREE.Scene();
    const geometry = new THREE.PlaneGeometry(2, 2);
    const uniforms = {
      time: { value: 1.0 },
      resolution: { value: new THREE.Vector2() },
    };
    const material = new THREE.ShaderMaterial({ uniforms, vertexShader, fragmentShader });
    scene.add(new THREE.Mesh(geometry, material));
    const renderer = new THREE.WebGLRenderer({ alpha: true });
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
    container.appendChild(renderer.domElement);

    const resize = () => {
      renderer.setSize(Math.max(1, container.clientWidth), Math.max(1, container.clientHeight), false);
      uniforms.resolution.value.set(renderer.domElement.width, renderer.domElement.height);
    };
    const animate = () => {
      const animationId = requestAnimationFrame(animate);
      uniforms.time.value += 0.05;
      renderer.render(scene, camera);
      if (sceneRef.current) sceneRef.current.animationId = animationId;
    };

    sceneRef.current = { camera, scene, renderer, uniforms, animationId: 0, geometry, material };
    resize();
    window.addEventListener("resize", resize);
    animate();

    return () => {
      window.removeEventListener("resize", resize);
      if (sceneRef.current) {
        cancelAnimationFrame(sceneRef.current.animationId);
        if (sceneRef.current.renderer.domElement.parentElement === container) {
          container.removeChild(sceneRef.current.renderer.domElement);
        }
        sceneRef.current.renderer.dispose();
        sceneRef.current.geometry.dispose();
        sceneRef.current.material.dispose();
      }
      sceneRef.current = null;
    };
  }, []);

  return <div ref={containerRef} className={`shader-canvas shader-canvas--lines ${className}`} />;
}

export const ShaderAnimation = ShaderLines;
