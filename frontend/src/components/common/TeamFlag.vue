<script setup lang="ts">
import { computed, ref, watch } from 'vue'

const teamIconModules = import.meta.glob('../../assets/team-icons/*.png', {
  eager: true,
  import: 'default',
}) as Record<string, string>

const teamIcons = Object.fromEntries(
  Object.entries(teamIconModules).map(([path, url]) => {
    const filename = path.split('/').pop() || ''
    return [filename.replace(/\.png$/i, '').toUpperCase(), url]
  }),
)

const props = withDefaults(defineProps<{
  value?: string | null
  alt?: string
  fallback?: string
  size?: 'sm' | 'md' | 'lg'
}>(), {
  value: '',
  alt: '',
  fallback: '',
  size: 'md',
})

const failed = ref(false)

watch(() => [props.value, props.fallback], () => {
  failed.value = false
})

const source = computed(() => (props.value || '').trim())
const fallbackText = computed(() => {
  const text = (props.fallback || props.alt || '').trim()
  return text ? text.slice(0, 3).toUpperCase() : 'TBD'
})
const code = computed(() => {
  const fallback = (props.fallback || '').trim().toUpperCase()
  if (/^[A-Z]{2,4}$/.test(fallback)) return fallback
  const sourceText = source.value.trim().toUpperCase()
  if (/^[A-Z]{2,4}$/.test(sourceText)) return sourceText
  return ''
})
const artCodes = new Set([
  'ALB', 'ARG', 'AUS', 'BEL', 'BRA', 'CAN', 'COL', 'CPV', 'CRC', 'CRO',
  'CZE', 'EGY', 'ENG', 'ESP', 'FRA', 'FRO', 'GER', 'GHA', 'IRN', 'ITA',
  'JPN', 'KAZ', 'KOR', 'KSA', 'LVA', 'MAR', 'MEX', 'MKD', 'MNE', 'NED',
  'NOR', 'NZL', 'PAR', 'POR', 'QAT', 'RSA', 'SEN', 'SRB', 'SUI', 'TUN',
  'URU', 'USA', 'WAL',
])
const artClass = computed(() => artCodes.has(code.value) ? `flag-${code.value.toLowerCase()}` : '')
const isImage = computed(() => /^(https?:\/\/|\/)/i.test(source.value))
const localIcon = computed(() => code.value ? teamIcons[code.value] : '')
const imageSrc = computed(() => localIcon.value || (isImage.value ? source.value : ''))
const isLocalIcon = computed(() => Boolean(localIcon.value && !failed.value))
</script>

<template>
  <span class="team-flag" :class="[size, { 'has-local-icon': isLocalIcon }]" aria-hidden="true">
    <img
      v-if="imageSrc && !failed"
      :src="imageSrc"
      :alt="alt"
      loading="lazy"
      decoding="async"
      @error="failed = true"
    />
    <span v-else-if="artClass" class="flag-art" :class="artClass"></span>
    <span v-else-if="source && !isImage" class="flag-text">{{ source }}</span>
    <span v-else class="flag-fallback">{{ fallbackText }}</span>
  </span>
</template>

<style scoped>
.team-flag {
  --flag-size: 34px;
  width: var(--flag-size);
  height: var(--flag-size);
  min-width: var(--flag-size);
  aspect-ratio: 1 / 1;
  display: inline-grid;
  place-items: center;
  flex: 0 0 auto;
  padding: 2px;
  overflow: hidden;
  border-radius: 999px;
  border: 1px solid rgba(15, 23, 42, 0.22);
  background: linear-gradient(145deg, #fff, #eef2f7);
  box-shadow:
    0 1px 2px rgba(15, 23, 42, 0.2),
    inset 0 1px 0 rgba(255, 255, 255, 0.95);
  color: var(--text);
  vertical-align: middle;
  box-sizing: border-box;
}

.team-flag.sm {
  --flag-size: 24px;
}

.team-flag.lg {
  --flag-size: 54px;
}

.team-flag.has-local-icon {
  padding: 0;
  border: 0;
  background: transparent;
  box-shadow: none;
}

.team-flag img {
  width: 100%;
  height: 100%;
  display: block;
  object-fit: contain;
  border-radius: inherit;
}

.team-flag.has-local-icon img {
  object-position: center;
}

.flag-art {
  position: relative;
  width: 100%;
  height: 100%;
  display: block;
  overflow: hidden;
  border-radius: inherit;
  background: #f8fafc;
}

.flag-art::before,
.flag-art::after {
  content: "";
  position: absolute;
  box-sizing: border-box;
}

.flag-mex,
.flag-ita {
  background: linear-gradient(90deg, #00843d 0 33.33%, #fff 33.33% 66.66%, #ce1126 66.66%);
}

.flag-mex::after {
  left: 43%;
  top: 42%;
  width: 14%;
  height: 16%;
  border-radius: 50%;
  background: #b8872f;
}

.flag-can {
  background: linear-gradient(90deg, #d80621 0 24%, #fff 24% 76%, #d80621 76%);
}

.flag-can::after {
  left: 42%;
  top: 31%;
  width: 16%;
  height: 30%;
  background: #d80621;
  clip-path: polygon(50% 0, 62% 30%, 86% 20%, 70% 48%, 100% 56%, 64% 62%, 70% 100%, 50% 75%, 30% 100%, 36% 62%, 0 56%, 30% 48%, 14% 20%, 38% 30%);
}

.flag-usa {
  background: repeating-linear-gradient(180deg, #b22234 0 7.7%, #fff 7.7% 15.4%);
}

.flag-usa::before {
  left: 0;
  top: 0;
  width: 54%;
  height: 54%;
  background:
    radial-gradient(circle at 22% 25%, #fff 0 4%, transparent 4.5%),
    radial-gradient(circle at 46% 42%, #fff 0 4%, transparent 4.5%),
    radial-gradient(circle at 70% 25%, #fff 0 4%, transparent 4.5%),
    radial-gradient(circle at 26% 68%, #fff 0 4%, transparent 4.5%),
    radial-gradient(circle at 68% 70%, #fff 0 4%, transparent 4.5%),
    #3c3b6e;
}

.flag-nzl,
.flag-aus {
  background:
    radial-gradient(circle at 74% 28%, #fff 0 3.4%, transparent 3.8%),
    radial-gradient(circle at 84% 47%, #fff 0 3.4%, transparent 3.8%),
    radial-gradient(circle at 67% 66%, #fff 0 3.4%, transparent 3.8%),
    radial-gradient(circle at 86% 75%, #fff 0 3.4%, transparent 3.8%),
    #00247d;
}

.flag-nzl::before,
.flag-aus::before {
  left: 0;
  top: 0;
  width: 44%;
  height: 36%;
  background:
    linear-gradient(90deg, transparent 42%, #fff 42% 58%, transparent 58%),
    linear-gradient(180deg, transparent 38%, #fff 38% 62%, transparent 62%),
    linear-gradient(90deg, transparent 46%, #cf142b 46% 54%, transparent 54%),
    linear-gradient(180deg, transparent 43%, #cf142b 43% 57%, transparent 57%),
    #012169;
}

.flag-nzl::after,
.flag-aus::after {
  display: none;
}

.flag-jpn {
  background: #fff;
}

.flag-jpn::after {
  left: 31%;
  top: 31%;
  width: 38%;
  height: 38%;
  border-radius: 50%;
  background: #bc002d;
}

.flag-irn {
  background: linear-gradient(180deg, #239f40 0 33.33%, #fff 33.33% 66.66%, #da0000 66.66%);
}

.flag-irn::after {
  left: 43%;
  top: 39%;
  width: 14%;
  height: 22%;
  border-radius: 50%;
  border: 2px solid #da0000;
}

.flag-esp {
  background: linear-gradient(180deg, #aa151b 0 25%, #f1bf00 25% 75%, #aa151b 75%);
}

.flag-esp::after,
.flag-crc::after,
.flag-egy::after,
.flag-srb::after,
.flag-cro::after {
  left: 27%;
  top: 43%;
  width: 11%;
  height: 14%;
  border-radius: 3px;
  background: #d6a530;
}

.flag-arg {
  background: linear-gradient(180deg, #74acdf 0 32%, #fff 32% 68%, #74acdf 68%);
}

.flag-arg::after,
.flag-uru::before {
  left: 43%;
  top: 43%;
  width: 14%;
  height: 14%;
  border-radius: 50%;
  background: #f6b40e;
}

.flag-mar {
  background: #c1272d;
}

.flag-mar::after {
  left: 35%;
  top: 35%;
  width: 30%;
  height: 30%;
  clip-path: polygon(50% 0, 61% 36%, 98% 36%, 68% 57%, 79% 92%, 50% 70%, 21% 92%, 32% 57%, 2% 36%, 39% 36%);
  background: #006233;
}

.flag-bra {
  background: #009b3a;
}

.flag-bra::before {
  left: 16%;
  top: 24%;
  width: 68%;
  height: 52%;
  background: #ffdf00;
  clip-path: polygon(50% 0, 100% 50%, 50% 100%, 0 50%);
}

.flag-bra::after {
  left: 35%;
  top: 35%;
  width: 30%;
  height: 30%;
  border-radius: 50%;
  background: #002776;
}

.flag-sui {
  background: #d52b1e;
}

.flag-sui::before {
  left: 27%;
  top: 42%;
  width: 46%;
  height: 16%;
  background: #fff;
}

.flag-sui::after {
  left: 42%;
  top: 27%;
  width: 16%;
  height: 46%;
  background: #fff;
}

.flag-par,
.flag-ned,
.flag-lva {
  background: linear-gradient(180deg, #ae1c28 0 33.33%, #fff 33.33% 66.66%, #21468b 66.66%);
}

.flag-lva {
  background: linear-gradient(180deg, #9e3039 0 42%, #fff 42% 58%, #9e3039 58%);
}

.flag-tun {
  background: #e70013;
}

.flag-tun::after {
  left: 30%;
  top: 30%;
  width: 40%;
  height: 40%;
  border-radius: 50%;
  background: radial-gradient(circle at 58% 50%, #e70013 0 23%, transparent 24%), #fff;
}

.flag-crc {
  background: linear-gradient(180deg, #002b7f 0 16%, #fff 16% 32%, #ce1126 32% 68%, #fff 68% 84%, #002b7f 84%);
}

.flag-ger {
  background: linear-gradient(180deg, #000 0 33.33%, #dd0000 33.33% 66.66%, #ffce00 66.66%);
}

.flag-kor {
  background: #fff;
}

.flag-kor::after {
  left: 34%;
  top: 34%;
  width: 32%;
  height: 32%;
  border-radius: 50%;
  background: conic-gradient(#cd2e3a 0 50%, #0047a0 0);
  transform: rotate(35deg);
}

.flag-ksa {
  background: #006c35;
}

.flag-ksa::after {
  left: 22%;
  bottom: 28%;
  width: 56%;
  height: 5%;
  border-radius: 999px;
  background: #fff;
}

.flag-bel {
  background: linear-gradient(90deg, #000 0 33.33%, #fae042 33.33% 66.66%, #ed2939 66.66%);
}

.flag-egy {
  background: linear-gradient(180deg, #ce1126 0 33.33%, #fff 33.33% 66.66%, #000 66.66%);
}

.flag-col {
  background: linear-gradient(180deg, #fcd116 0 50%, #003893 50% 75%, #ce1126 75%);
}

.flag-cpv {
  background: linear-gradient(180deg, #003893 0 48%, #fff 48% 54%, #cf2027 54% 60%, #fff 60% 66%, #003893 66%);
}

.flag-fra {
  background: linear-gradient(90deg, #002395 0 33.33%, #fff 33.33% 66.66%, #ed2939 66.66%);
}

.flag-sen {
  background: linear-gradient(90deg, #00853f 0 33.33%, #fdef42 33.33% 66.66%, #e31b23 66.66%);
}

.flag-sen::after {
  left: 43%;
  top: 36%;
  width: 14%;
  height: 20%;
  clip-path: polygon(50% 0, 61% 36%, 98% 36%, 68% 57%, 79% 92%, 50% 70%, 21% 92%, 32% 57%, 2% 36%, 39% 36%);
  background: #00853f;
}

.flag-uru {
  background: repeating-linear-gradient(180deg, #fff 0 11.11%, #0038a8 11.11% 22.22%);
}

.flag-uru::before {
  left: 0;
  top: 0;
}

.flag-qat {
  background: linear-gradient(90deg, #fff 0 30%, #8d1b3d 30%);
}

.flag-qat::after {
  left: 24%;
  top: 0;
  width: 18%;
  height: 100%;
  background: repeating-linear-gradient(180deg, transparent 0 9%, #8d1b3d 9% 18%);
}

.flag-nor,
.flag-fro {
  background:
    linear-gradient(90deg, transparent 31%, #fff 31% 45%, #00205b 45% 55%, #fff 55% 69%, transparent 69%),
    linear-gradient(180deg, transparent 31%, #fff 31% 45%, #00205b 45% 55%, #fff 55% 69%, transparent 69%),
    #ba0c2f;
}

.flag-fro {
  background:
    linear-gradient(90deg, transparent 31%, #0065bd 31% 39%, #ed2939 39% 49%, #0065bd 49% 57%, transparent 57%),
    linear-gradient(180deg, transparent 36%, #0065bd 36% 43%, #ed2939 43% 53%, #0065bd 53% 60%, transparent 60%),
    #fff;
}

.flag-wal {
  background: linear-gradient(180deg, #fff 0 50%, #00a650 50%);
}

.flag-wal::after {
  left: 28%;
  top: 33%;
  width: 44%;
  height: 26%;
  background: #d21034;
  clip-path: polygon(0 70%, 26% 36%, 44% 52%, 63% 20%, 100% 45%, 70% 78%, 30% 84%);
}

.flag-mkd {
  background: repeating-conic-gradient(from 0deg, #f9d616 0 12deg, #d91f26 12deg 30deg);
}

.flag-mkd::after {
  left: 35%;
  top: 35%;
  width: 30%;
  height: 30%;
  border-radius: 50%;
  background: #f9d616;
}

.flag-kaz {
  background: #00abc2;
}

.flag-kaz::after {
  left: 44%;
  top: 31%;
  width: 24%;
  height: 24%;
  border-radius: 50%;
  background: #f6c400;
}

.flag-eng {
  background:
    linear-gradient(90deg, transparent 42%, #ce1124 42% 58%, transparent 58%),
    linear-gradient(180deg, transparent 42%, #ce1124 42% 58%, transparent 58%),
    #fff;
}

.flag-srb {
  background: linear-gradient(180deg, #c6363c 0 33.33%, #0c4076 33.33% 66.66%, #fff 66.66%);
}

.flag-alb {
  background: #e41e20;
}

.flag-alb::after {
  left: 34%;
  top: 30%;
  width: 32%;
  height: 36%;
  border-radius: 50%;
  background: #111;
}

.flag-cro {
  background: linear-gradient(180deg, #f00 0 33.33%, #fff 33.33% 66.66%, #171796 66.66%);
}

.flag-cze {
  background: linear-gradient(180deg, #fff 0 50%, #d7141a 50%);
}

.flag-cze::before {
  left: 0;
  top: 0;
  width: 55%;
  height: 100%;
  background: #11457e;
  clip-path: polygon(0 0, 100% 50%, 0 100%);
}

.flag-mne {
  background: #c40308;
  border: 3px solid #d4af37;
}

.flag-mne::after {
  left: 34%;
  top: 30%;
  width: 32%;
  height: 36%;
  border-radius: 50%;
  background: #d4af37;
}

.flag-por {
  background: linear-gradient(90deg, #006600 0 40%, #ff0000 40%);
}

.flag-por::after {
  left: 34%;
  top: 38%;
  width: 17%;
  height: 24%;
  border-radius: 50%;
  background: #ffcc00;
}

.flag-gha {
  background: linear-gradient(180deg, #ce1126 0 33.33%, #fcd116 33.33% 66.66%, #006b3f 66.66%);
}

.flag-gha::after {
  left: 43%;
  top: 41%;
  width: 14%;
  height: 18%;
  clip-path: polygon(50% 0, 61% 36%, 98% 36%, 68% 57%, 79% 92%, 50% 70%, 21% 92%, 32% 57%, 2% 36%, 39% 36%);
  background: #111;
}

.flag-rsa {
  background:
    linear-gradient(145deg, transparent 0 39%, #007a4d 39% 52%, transparent 52%),
    linear-gradient(35deg, #000 0 28%, #ffb612 28% 34%, #007a4d 34% 48%, transparent 48%),
    linear-gradient(180deg, #de3831 0 50%, #002395 50%);
}

.flag-text {
  max-width: 100%;
  overflow: hidden;
  font-size: calc(var(--flag-size) * 0.62);
  line-height: 1;
  text-align: center;
}

.flag-fallback {
  max-width: 100%;
  padding: 0 3px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: var(--muted);
  font-size: calc(var(--flag-size) * 0.28);
  font-weight: 800;
  letter-spacing: 0;
}
</style>
