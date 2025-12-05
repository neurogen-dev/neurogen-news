/**
 * ═══════════════════════════════════════════════════════════════════════════
 * VELVET GLASS DESIGN SYSTEM - Interactive Components
 * ═══════════════════════════════════════════════════════════════════════════
 */

(function() {
  'use strict';

  // ─────────────────────────────────────────────────────────────────────────
  // UTILITIES
  // ─────────────────────────────────────────────────────────────────────────
  
  const $ = (selector, context = document) => context.querySelector(selector);
  const $$ = (selector, context = document) => [...context.querySelectorAll(selector)];
  
  const debounce = (fn, delay) => {
    let timeoutId;
    return (...args) => {
      clearTimeout(timeoutId);
      timeoutId = setTimeout(() => fn.apply(this, args), delay);
    };
  };

  // ─────────────────────────────────────────────────────────────────────────
  // RIPPLE EFFECT
  // ─────────────────────────────────────────────────────────────────────────
  
  class RippleEffect {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.btn-velvet, .btn-icon').forEach(btn => {
        btn.addEventListener('click', this.createRipple.bind(this));
      });
    }
    
    createRipple(e) {
      const btn = e.currentTarget;
      const rect = btn.getBoundingClientRect();
      
      const ripple = document.createElement('span');
      ripple.className = 'velvet-ripple';
      
      const size = Math.max(rect.width, rect.height) * 2;
      const x = e.clientX - rect.left - size / 2;
      const y = e.clientY - rect.top - size / 2;
      
      ripple.style.cssText = `
        position: absolute;
        width: ${size}px;
        height: ${size}px;
        left: ${x}px;
        top: ${y}px;
        background: radial-gradient(circle, rgba(255,255,255,0.4) 0%, transparent 70%);
        border-radius: 50%;
        transform: scale(0);
        animation: rippleExpand 0.6s cubic-bezier(0.4, 0, 0.2, 1) forwards;
        pointer-events: none;
      `;
      
      btn.appendChild(ripple);
      
      setTimeout(() => ripple.remove(), 600);
    }
  }
  
  // Add ripple keyframes
  const style = document.createElement('style');
  style.textContent = `
    @keyframes rippleExpand {
      0% {
        transform: scale(0);
        opacity: 1;
      }
      100% {
        transform: scale(1);
        opacity: 0;
      }
    }
  `;
  document.head.appendChild(style);

  // ─────────────────────────────────────────────────────────────────────────
  // MAGNETIC HOVER EFFECT
  // ─────────────────────────────────────────────────────────────────────────
  
  class MagneticEffect {
    constructor() {
      this.strength = 0.3;
      this.init();
    }
    
    init() {
      $$('.btn-icon, .btn-social').forEach(el => {
        el.addEventListener('mousemove', this.handleMouseMove.bind(this));
        el.addEventListener('mouseleave', this.handleMouseLeave.bind(this));
      });
    }
    
    handleMouseMove(e) {
      const el = e.currentTarget;
      const rect = el.getBoundingClientRect();
      const centerX = rect.left + rect.width / 2;
      const centerY = rect.top + rect.height / 2;
      
      const deltaX = (e.clientX - centerX) * this.strength;
      const deltaY = (e.clientY - centerY) * this.strength;
      
      el.style.transform = `translate(${deltaX}px, ${deltaY}px)`;
    }
    
    handleMouseLeave(e) {
      e.currentTarget.style.transform = '';
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // TABS COMPONENT
  // ─────────────────────────────────────────────────────────────────────────
  
  class VelvetTabs {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.tabs-velvet').forEach(tabGroup => {
        const tabs = $$('.tab-item', tabGroup);
        tabs.forEach((tab, index) => {
          tab.addEventListener('click', () => this.switchTab(tabs, index));
        });
      });
    }
    
    switchTab(tabs, activeIndex) {
      tabs.forEach((tab, index) => {
        if (index === activeIndex) {
          tab.classList.add('active');
        } else {
          tab.classList.remove('active');
        }
      });
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // CHIPS COMPONENT
  // ─────────────────────────────────────────────────────────────────────────
  
  class VelvetChips {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.chips-group').forEach(group => {
        const chips = $$('.chip-velvet', group);
        chips.forEach((chip, index) => {
          chip.addEventListener('click', () => this.toggleChip(chips, index));
        });
      });
    }
    
    toggleChip(chips, index) {
      const chip = chips[index];
      
      // If this is the "All" chip (first one with dot), make it exclusive
      if (chip.querySelector('.chip-dot')) {
        chips.forEach(c => c.classList.remove('active'));
        chip.classList.add('active');
      } else {
        // Remove "All" selection when selecting specific chips
        const allChip = chips.find(c => c.querySelector('.chip-dot'));
        if (allChip) allChip.classList.remove('active');
        
        chip.classList.toggle('active');
        
        // If no chips selected, select "All"
        const anyActive = chips.some(c => c.classList.contains('active'));
        if (!anyActive && allChip) {
          allChip.classList.add('active');
        }
      }
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // SLIDER COMPONENT
  // ─────────────────────────────────────────────────────────────────────────
  
  class VelvetSlider {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.slider-velvet').forEach(slider => {
        const input = $('input[type="range"]', slider);
        const fill = $('.slider-fill', slider);
        const thumb = $('.slider-thumb', slider);
        const value = $('.slider-value', slider);
        
        if (input) {
          this.updateSlider(input, fill, thumb, value);
          
          input.addEventListener('input', () => {
            this.updateSlider(input, fill, thumb, value);
          });
        }
      });
    }
    
    updateSlider(input, fill, thumb, valueEl) {
      const percent = ((input.value - input.min) / (input.max - input.min)) * 100;
      
      if (fill) fill.style.width = `${percent}%`;
      if (thumb) thumb.style.left = `${percent}%`;
      if (valueEl) valueEl.textContent = `${Math.round(percent)}%`;
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // PROGRESS ANIMATION
  // ─────────────────────────────────────────────────────────────────────────
  
  class ProgressAnimation {
    constructor() {
      this.observed = new Set();
      this.init();
    }
    
    init() {
      const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting && !this.observed.has(entry.target)) {
            this.observed.add(entry.target);
            this.animateProgress(entry.target);
          }
        });
      }, { threshold: 0.5 });
      
      $$('.progress-fill').forEach(el => observer.observe(el));
    }
    
    animateProgress(el) {
      const width = el.style.width || '0%';
      el.style.width = '0%';
      
      requestAnimationFrame(() => {
        el.style.transition = 'width 1s cubic-bezier(0.4, 0, 0.2, 1)';
        el.style.width = width;
      });
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // CLOCK WIDGET
  // ─────────────────────────────────────────────────────────────────────────
  
  class ClockWidget {
    constructor() {
      this.init();
    }
    
    init() {
      this.updateClock();
      setInterval(() => this.updateClock(), 1000);
    }
    
    updateClock() {
      const now = new Date();
      
      // Analog hands
      const secondHand = $('.second-hand');
      const minuteHand = $('.minute-hand');
      const hourHand = $('.hour-hand');
      
      if (secondHand) {
        const seconds = now.getSeconds();
        const minutes = now.getMinutes();
        const hours = now.getHours() % 12;
        
        const secondDeg = seconds * 6;
        const minuteDeg = minutes * 6 + seconds * 0.1;
        const hourDeg = hours * 30 + minutes * 0.5;
        
        secondHand.style.transform = `translateX(-50%) rotate(${secondDeg}deg)`;
        minuteHand.style.transform = `translateX(-50%) rotate(${minuteDeg}deg)`;
        hourHand.style.transform = `translateX(-50%) rotate(${hourDeg}deg)`;
      }
      
      // Digital display
      const digitalTime = $('.digital-time');
      const digitalDate = $('.digital-date');
      
      if (digitalTime) {
        digitalTime.textContent = now.toLocaleTimeString('ru-RU', {
          hour: '2-digit',
          minute: '2-digit'
        });
      }
      
      if (digitalDate) {
        digitalDate.textContent = now.toLocaleDateString('ru-RU', {
          weekday: 'long',
          day: 'numeric',
          month: 'long'
        });
      }
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // TOAST NOTIFICATIONS
  // ─────────────────────────────────────────────────────────────────────────
  
  class ToastManager {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.toast-close').forEach(btn => {
        btn.addEventListener('click', (e) => {
          const toast = e.currentTarget.closest('.toast-velvet');
          this.dismiss(toast);
        });
      });
    }
    
    dismiss(toast) {
      toast.style.transition = 'all 0.3s cubic-bezier(0.4, 0, 0.2, 1)';
      toast.style.transform = 'translateX(20px)';
      toast.style.opacity = '0';
      
      setTimeout(() => {
        toast.style.display = 'none';
      }, 300);
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // CARD MOUSE TRACKING
  // ─────────────────────────────────────────────────────────────────────────
  
  class CardMouseTracker {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.card-glass, .philosophy-card').forEach(card => {
        card.addEventListener('mousemove', this.handleMouseMove.bind(this));
        card.addEventListener('mouseleave', this.handleMouseLeave.bind(this));
      });
    }
    
    handleMouseMove(e) {
      const card = e.currentTarget;
      const rect = card.getBoundingClientRect();
      
      const x = ((e.clientX - rect.left) / rect.width) * 100;
      const y = ((e.clientY - rect.top) / rect.height) * 100;
      
      card.style.setProperty('--mouse-x', `${x}%`);
      card.style.setProperty('--mouse-y', `${y}%`);
      
      // Add dynamic lighting effect
      const glowX = x - 50;
      const glowY = y - 50;
      card.style.background = `
        radial-gradient(
          ellipse at ${x}% ${y}%,
          rgba(255, 255, 255, 0.15) 0%,
          transparent 50%
        ),
        rgba(255, 255, 255, 0.72)
      `;
    }
    
    handleMouseLeave(e) {
      e.currentTarget.style.background = '';
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // SMOOTH SCROLL ANIMATIONS
  // ─────────────────────────────────────────────────────────────────────────
  
  class ScrollAnimations {
    constructor() {
      this.init();
    }
    
    init() {
      const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
          if (entry.isIntersecting) {
            entry.target.classList.add('is-visible');
          }
        });
      }, {
        threshold: 0.1,
        rootMargin: '0px 0px -50px 0px'
      });
      
      $$('.section, .component-group, .philosophy-card, .card-velvet').forEach(el => {
        el.classList.add('scroll-animate');
        observer.observe(el);
      });
    }
  }
  
  // Add scroll animation styles
  const scrollStyle = document.createElement('style');
  scrollStyle.textContent = `
    .scroll-animate {
      opacity: 0;
      transform: translateY(30px);
      transition: opacity 0.6s cubic-bezier(0.4, 0, 0.2, 1), 
                  transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);
    }
    
    .scroll-animate.is-visible {
      opacity: 1;
      transform: translateY(0);
    }
    
    .scroll-animate:nth-child(2) { transition-delay: 0.1s; }
    .scroll-animate:nth-child(3) { transition-delay: 0.2s; }
    .scroll-animate:nth-child(4) { transition-delay: 0.3s; }
  `;
  document.head.appendChild(scrollStyle);

  // ─────────────────────────────────────────────────────────────────────────
  // PARALLAX ORB MOVEMENT
  // ─────────────────────────────────────────────────────────────────────────
  
  class ParallaxOrbs {
    constructor() {
      this.orbs = $$('.orb');
      this.init();
    }
    
    init() {
      window.addEventListener('mousemove', debounce((e) => {
        this.moveOrbs(e);
      }, 10));
    }
    
    moveOrbs(e) {
      const { clientX, clientY } = e;
      const centerX = window.innerWidth / 2;
      const centerY = window.innerHeight / 2;
      
      this.orbs.forEach((orb, index) => {
        const speed = (index + 1) * 0.02;
        const x = (clientX - centerX) * speed;
        const y = (clientY - centerY) * speed;
        
        orb.style.transform = `translate(${x}px, ${y}px)`;
      });
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // PASSWORD TOGGLE
  // ─────────────────────────────────────────────────────────────────────────
  
  class PasswordToggle {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.input-action').forEach(btn => {
        btn.addEventListener('click', (e) => {
          const wrapper = e.currentTarget.closest('.input-velvet');
          const input = wrapper.querySelector('input');
          
          if (input.type === 'password') {
            input.type = 'text';
            btn.innerHTML = `
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            `;
          } else {
            input.type = 'password';
            btn.innerHTML = `
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
            `;
          }
        });
      });
    }
  }

  // ─────────────────────────────────────────────────────────────────────────
  // COLOR SWATCH COPY
  // ─────────────────────────────────────────────────────────────────────────
  
  class SwatchCopy {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.color-swatch').forEach(swatch => {
        swatch.style.cursor = 'pointer';
        swatch.title = 'Нажмите, чтобы скопировать цвет';
        
        swatch.addEventListener('click', () => {
          const value = swatch.querySelector('.swatch-value')?.textContent;
          if (value) {
            navigator.clipboard.writeText(value).then(() => {
              this.showCopied(swatch);
            });
          }
        });
      });
    }
    
    showCopied(swatch) {
      const feedback = document.createElement('span');
      feedback.textContent = 'Скопировано!';
      feedback.style.cssText = `
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        background: rgba(0, 0, 0, 0.8);
        color: white;
        padding: 8px 16px;
        border-radius: 8px;
        font-size: 12px;
        font-weight: 500;
        animation: copiedFade 1.5s forwards;
        z-index: 10;
      `;
      
      swatch.style.position = 'relative';
      swatch.appendChild(feedback);
      
      setTimeout(() => feedback.remove(), 1500);
    }
  }
  
  // Add copied animation
  const copiedStyle = document.createElement('style');
  copiedStyle.textContent = `
    @keyframes copiedFade {
      0% { opacity: 0; transform: translate(-50%, -50%) scale(0.8); }
      20% { opacity: 1; transform: translate(-50%, -50%) scale(1); }
      80% { opacity: 1; }
      100% { opacity: 0; transform: translate(-50%, -60%); }
    }
  `;
  document.head.appendChild(copiedStyle);

  // ─────────────────────────────────────────────────────────────────────────
  // LIKE BUTTON ANIMATION
  // ─────────────────────────────────────────────────────────────────────────
  
  class LikeButton {
    constructor() {
      this.init();
    }
    
    init() {
      $$('.player-like, .btn-icon-raised').forEach(btn => {
        btn.addEventListener('click', () => this.toggle(btn));
      });
    }
    
    toggle(btn) {
      const isLiked = btn.classList.toggle('is-liked');
      const svg = btn.querySelector('svg');
      
      if (isLiked) {
        svg.setAttribute('fill', 'currentColor');
        btn.style.color = 'var(--color-danger)';
        
        // Heart burst animation
        this.createBurst(btn);
      } else {
        svg.setAttribute('fill', 'none');
        btn.style.color = '';
      }
    }
    
    createBurst(btn) {
      for (let i = 0; i < 6; i++) {
        const particle = document.createElement('span');
        const angle = (i / 6) * 360;
        const distance = 30 + Math.random() * 20;
        
        particle.style.cssText = `
          position: absolute;
          width: 6px;
          height: 6px;
          background: var(--color-danger);
          border-radius: 50%;
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          animation: burst 0.5s cubic-bezier(0.4, 0, 0.2, 1) forwards;
          --angle: ${angle}deg;
          --distance: ${distance}px;
        `;
        
        btn.style.position = 'relative';
        btn.appendChild(particle);
        
        setTimeout(() => particle.remove(), 500);
      }
    }
  }
  
  // Add burst animation
  const burstStyle = document.createElement('style');
  burstStyle.textContent = `
    @keyframes burst {
      0% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 1;
      }
      100% {
        transform: 
          translate(-50%, -50%) 
          rotate(var(--angle)) 
          translateY(calc(var(--distance) * -1))
          scale(0);
        opacity: 0;
      }
    }
  `;
  document.head.appendChild(burstStyle);

  // ─────────────────────────────────────────────────────────────────────────
  // INITIALIZE ALL
  // ─────────────────────────────────────────────────────────────────────────
  
  document.addEventListener('DOMContentLoaded', () => {
    new RippleEffect();
    new MagneticEffect();
    new VelvetTabs();
    new VelvetChips();
    new VelvetSlider();
    new ProgressAnimation();
    new ClockWidget();
    new ToastManager();
    new CardMouseTracker();
    new ScrollAnimations();
    new ParallaxOrbs();
    new PasswordToggle();
    new SwatchCopy();
    new LikeButton();
    
    console.log('✨ Velvet Glass Design System initialized');
  });
  
})();



