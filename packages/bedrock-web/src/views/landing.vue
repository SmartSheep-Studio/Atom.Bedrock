<template>
  <div class="container">
    <div class="text-center pt-36">
      <img src="../assets/icon.png" width="128" height="128" />
      <div class="mt-4">
        <div class="text-4xl font-bold">{{ $t("brand.name") }}</div>
        <div class="text-lg">{{ $t("brand.slogan") }}</div>
      </div>
    </div>

    <n-grid class="mt-16" item-responsive responsive="screen">
      <n-gi span="24 m:10" class="text-center">
        <img src="../assets/undraw/undraw_work_time_re_hdyv.svg" width="400" height="300" />
      </n-gi>
      <n-gi span="24 m:14 l:12" class="flex items-center sm:text-center md:text-left">
        <div>
          <div class="text-xl font-bold">{{ $t("pages.landing.feat.does-work") }}</div>
          <div class="text-lg mt-2">
            {{ $t("pages.landing.feat.does-work.desc") }}
          </div>
        </div>
      </n-gi>
    </n-grid>
  </div>
</template>
