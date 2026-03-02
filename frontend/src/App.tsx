import type { Component } from 'solid-js';
import { createSignal, createResource, For } from 'solid-js';
import TimePicker from "@rnwonder/solid-date-picker/timePicker";
import DatePicker, {
  PickerValue,
  TimeValue,
} from '@rnwonder/solid-date-picker';

import { CalendarIcon } from './Icons/Calendar';
import { NextIcon } from './Icons/Next';
import { PrevIcon } from './Icons/Prev';

import AgGridSolid from 'ag-grid-solid';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';
import { GridOptions } from 'ag-grid-community';

// const apiURL = 'http://vwsstorage.localdomain:3000/api/v1';
// @ts-ignore
const apiURL = window.base_url;
const apiGetDevices = apiURL + '/device';
const apiGetDevice = apiURL + '/device/{id}';
const apiGetDeviceEvents = apiURL + '/device/{id}/event?start={start}&end={end}';

interface IDevice {
  Id: string,
  name: string,
  label: string,
  type: string,
  room: string
}

const Device = {
  Id: '',
  name: '',
  label: '',
  type: '',
  room: ''
}

async function fetchAPI(apiUri: string) {
  console.log('fetchAPI: ', apiUri)
  try {
    const response = (await fetch(apiUri)).json();
    return response
  } catch (error) {
    console.error('Error fetching data:', error);
    throw error
  }
}

const App: Component = () => {
  const [area, setArea] = createSignal('device_details')
  const [readDevice, setReadDevice] = createSignal<number>(1);
  const toggleDevice = () => setReadDevice(readDevice() == 1 ? setReadDevice(2) : setReadDevice(1))
  const [devices] = createResource(readDevice, async () => await fetchAPI(apiGetDevices));

  const [readDeviceDetail, setReadDeviceDetail] = createSignal<IDevice>(Device);
  const [deviceDetail] = createResource(readDeviceDetail, async () => await fetchAPI(apiGetDevice.replace('{id}', readDeviceDetail().Id)));

  const [toggleReadDeviceEvents, settoggleReadDeviceEvents] = createSignal<number>(1);
  const toggleDeviceEvents = () => toggleReadDeviceEvents() == 1 ? settoggleReadDeviceEvents(2) : settoggleReadDeviceEvents(1)
  const [deviceEventsToggle] = createResource(toggleReadDeviceEvents, async () => await fetchAPI(apiGetDeviceEvents.replace('{id}', readDeviceDetail().Id).replace('{start}', buildDateTime('from')).replace('{end}', buildDateTime('to'))));

  const [singleDateFrom, setSingleDateFrom] = createSignal<PickerValue>({
    label: (new Date(buildDateTime("initFromDate")).toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("initFromDate")
    },
  });
  const [ldeDF] = createResource(singleDateFrom, loadDeviceEvents)

  const [timeValueFrom, setTimeValueFrom] = createSignal<TimeValue>({
    value: {
      hour: parseInt(buildDateTime("initFromHour")),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date(buildDateTime("initFromDate")).toLocaleTimeString()),
  });
  const [ldeTF] = createResource(timeValueFrom, loadDeviceEvents)

  const [singleDateTo, setSingleDateTo] = createSignal<PickerValue>({
    label: (new Date().toLocaleString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })),
    value: {
      selected: buildDateTime("current")
    },
  });
  const [ldeDT] = createResource(singleDateTo, loadDeviceEvents)

  const [timeValueTo, setTimeValueTo] = createSignal<TimeValue>({
    value: {
      hour: new Date().getHours(),
      minute: new Date().getMinutes(),
      second: new Date().getSeconds()
    },
    label: (new Date().toLocaleTimeString()),
  });
  const [ldeTT] = createResource(timeValueTo, loadDeviceEvents)
  const [deviceSelected, setDeviceSelected] = createSignal('false')

  function loadDevice(device: IDevice) {
    console.log('loadDevice: ', device)
    setDeviceSelected('true')
    console.log('deviceSelected: ', deviceSelected())
    setReadDeviceDetail(device)
    loadDeviceEvents()
  }

  function loadDeviceEvents() {
    console.log('loadDeviceEvents: ', readDeviceDetail)
    toggleDeviceEvents()
  }

  function buildDateTime(fromOrtoOrCurrent: string): string {

    function fixZeroDateTime(dateTime: string): string {
      let receivedDateTime = new Date(dateTime)
      let hours = receivedDateTime.getHours()
      let minutes = receivedDateTime.getMinutes()
      let seconds = receivedDateTime.getSeconds()

      let currenthours = ("0" + hours).slice(-2);
      let currentminutes = ("0" + minutes).slice(-2);
      let currentseconds = ("0" + seconds).slice(-2);

      let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE') + "T" + currenthours + ":" + currentminutes + ":" + currentseconds + ".000"
      //let returnFixedZeroDateTime = receivedDateTime.toLocaleDateString('sv-SE')+"T"+currenthours+":"+currentminutes+":"+currentseconds+".000Z"

      return returnFixedZeroDateTime
    }

    if (fromOrtoOrCurrent == 'from') {
      let fromDateTime = singleDateFrom().value?.selected?.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
      console.log('fromDateTime: ', fromDateTime)
      //let fromDateTime = singleDateFrom().value.selected.split('T')[0] + 'T' + timeValueFrom().value.hour + ':' + timeValueFrom().value.minute + ':' + timeValueFrom().value.second
      return fromDateTime
    } else if (fromOrtoOrCurrent == 'to') {
      let toDateTime = singleDateTo().value?.selected?.split('T')[0] + 'T' + timeValueTo().value.hour + ':' + timeValueTo().value.minute + ':' + timeValueTo().value.second
      console.log('toDateTime: ', toDateTime)
      return toDateTime
    } else if (fromOrtoOrCurrent == 'current') {
      let currentDateTime = new Date()
      let returnCurrentDateTime = fixZeroDateTime(currentDateTime.toLocaleString())
      console.log('returnCurrentDateTime: ', returnCurrentDateTime)
      return returnCurrentDateTime
    } else if (fromOrtoOrCurrent == 'initFromDate') {
      let initFromDateTime = new Date()
      initFromDateTime.setHours(initFromDateTime.getHours() - 8)

      let initFromDateTimeStr = fixZeroDateTime(initFromDateTime.toLocaleString())
      console.log('initFromDateTime: ', initFromDateTimeStr)
      return initFromDateTimeStr
    } else {
      let initFromHour = new Date()
      initFromHour.setHours(initFromHour.getHours() - 8)
      let returnInitFromHour = initFromHour.getHours()
      console.log('initFromHour: ', returnInitFromHour)
      return returnInitFromHour.toString()
    }
  }

  const defaultColDef = {
    flex: 1,
    //minWidth: 100,
    filter: true,
    sortable: true,
    resizable: true,
  };

  const gridOptions: GridOptions = {
    autoSizeStrategy: {
      type: 'fitCellContents',
      //defaultMinWidth: 100,

      // columnLimits: [
      //   {
      //     colId: 'descriptionText',
      //     minWidth: 900
      //   }
      // ]
    },

    // other grid options ...
  }

  const columnDefsDevice = [
    { field: 'Id' },
    { field: 'name' },
    { field: 'label' },
    { field: 'type' },
    //{ field: 'room' },
  ];

  const columnDefsEvent = [
    { field: 'TimeStamp' },
    { field: 'name' },
    { field: 'value' },
    { field: 'displayName' },
    { field: 'deviceId' },
    { field: 'descriptionText' },
    { field: 'unit' },
    { field: 'type' },
    { field: 'data' },
  ];

  toggleDevice()

  return (
    <div class="flex flex-col h-screen bg-slate-50">
      {/* ── Header ─────────────────────────────────────── */}
      <header class="bg-slate-900 text-white flex items-center px-5 h-14 shrink-0 shadow-lg">
        <svg class="w-6 h-6 text-indigo-400 mr-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
        </svg>
        <h1 class="text-lg font-semibold tracking-wide">Hubitat Device Events</h1>
      </header>

      <div class="flex flex-1 overflow-hidden">
        {/* ── Sidebar ────────────────────────────────────── */}
        <aside class="w-72 bg-slate-900 flex flex-col shrink-0 border-r border-slate-800">
          <div class="px-4 py-3 border-b border-slate-800">
            <p class="text-xs font-semibold uppercase tracking-wider text-slate-400">Devices</p>
          </div>
          <div class="flex-1 overflow-y-auto px-2 py-2 space-y-0.5">
            {devices.loading && (
              <p class="text-sm text-slate-500 px-3 py-4 text-center">Loading devices…</p>
            )}
            <For each={devices()}>
              {(device) => (
                <div
                  class={'device-item' + (readDeviceDetail().Id === device.Id ? ' active' : '')}
                  onClick={() => loadDevice(device)}
                >
                  {device.label}
                </div>
              )}
            </For>
          </div>
        </aside>

        {/* ── Main content ───────────────────────────────── */}
        <main class="flex-1 flex flex-col overflow-hidden">
          {/* Tabs */}
          <nav class="flex border-b border-slate-200 bg-white px-4 shrink-0">
            <button
              class={'tab-btn' + (area() == 'device_details' ? ' tab-active' : '')}
              onclick={() => setArea('device_details')}
            >
              Device Details
            </button>
            <button
              class={'tab-btn' + (area() == 'device_events' ? ' tab-active' : '')}
              onclick={() => setArea('device_events')}
              disabled={deviceSelected() == 'false'}
            >
              Device Events
            </button>
          </nav>

          {/* ── Device Details tab ─────────────────────── */}
          <div class={'flex-1 overflow-hidden ' + (area() == 'device_details' ? '' : 'hidden')}>
            {deviceSelected() == 'false' ? (
              <div class="flex items-center justify-center h-full text-slate-400">
                <div class="text-center">
                  <svg class="w-16 h-16 mx-auto mb-4 text-slate-300" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
                  </svg>
                  <p class="text-lg font-medium">Select a device</p>
                  <p class="text-sm mt-1">Choose a device from the sidebar to view its details</p>
                </div>
              </div>
            ) : (
              <div class="ag-theme-alpine h-full w-full">
                <AgGridSolid
                  rowData={deviceDetail()}
                  columnDefs={columnDefsDevice}
                  defaultColDef={defaultColDef}
                  animateRows={true}
                  autoSizeStrategy={gridOptions.autoSizeStrategy}
                  rowSelection='single'
                />
              </div>
            )}
          </div>

          {/* ── Device Events tab ──────────────────────── */}
          <div class={'flex-1 flex flex-col overflow-hidden ' + (area() == 'device_events' ? '' : 'hidden')}>
            {/* Date/time toolbar */}
            <div class="bg-white border-b border-slate-200 px-5 py-3 flex flex-wrap items-center gap-6 shrink-0">
              <div class="flex items-center gap-2">
                <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">From</span>
                <DatePicker
                  removeNavButtons
                  weekDaysType="single"
                  shouldHighlightWeekends
                  hideOutSideDays
                  afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                    <div class="flex gap-1">
                      <button onClick={handlePrevMonth}><PrevIcon /></button>
                      <button onClick={handleNextMonth}><NextIcon /></button>
                    </div>
                  )}
                  renderInput={({ showDate }) => (
                    <div class="custom-input">
                      <input value={singleDateFrom().label} readOnly />
                      <button onClick={showDate}><CalendarIcon /></button>
                    </div>
                  )}
                  onChange={() => { toggleDeviceEvents }}
                  value={singleDateFrom}
                  setValue={setSingleDateFrom}
                  type={'single'}
                  monthSelectorFormat={'long'}
                />
                <TimePicker
                  value={timeValueFrom}
                  setValue={setTimeValueFrom}
                  shouldCloseOnSelect
                  allowedView={['hour', 'minute', 'second']}
                  renderInput={({ showTime }) => (
                    <div class="custom-input">
                      <input value={timeValueFrom().label} readOnly />
                      <button onClick={showTime}><CalendarIcon /></button>
                    </div>
                  )}
                  timeAnalogClockCenterDotClass='text-indigo-500'
                  timeAnalogNumberClass='text-slate-700'
                  timeAnalogClockHandClass='text-indigo-500'
                  timeAnalogWrapperClass='text-slate-700'
                  timePickerBottomAreaClass='text-slate-700'
                  timePickerMeridiemBtnClass='text-slate-700'
                  timePickerTopAreaClass='text-slate-700'
                  timePickerWrapperClass='text-slate-700'
                />
              </div>

              <div class="w-px h-6 bg-slate-300"></div>

              <div class="flex items-center gap-2">
                <span class="text-xs font-semibold uppercase tracking-wider text-slate-400">To</span>
                <DatePicker
                  removeNavButtons
                  weekDaysType="single"
                  shouldHighlightWeekends
                  hideOutSideDays
                  afterNextButtonAreaJSX={({ handleNextMonth, handlePrevMonth }) => (
                    <div class="flex gap-1">
                      <button onClick={handlePrevMonth}><PrevIcon /></button>
                      <button onClick={handleNextMonth}><NextIcon /></button>
                    </div>
                  )}
                  renderInput={({ showDate }) => (
                    <div class="custom-input">
                      <input value={singleDateTo().label} readOnly />
                      <button onClick={showDate}><CalendarIcon /></button>
                    </div>
                  )}
                  onChange={() => { toggleDeviceEvents }}
                  value={singleDateTo}
                  setValue={setSingleDateTo}
                  type={'single'}
                  monthSelectorFormat={'long'}
                />
                <TimePicker
                  value={timeValueTo}
                  setValue={setTimeValueTo}
                  shouldCloseOnSelect
                  allowedView={['hour', 'minute', 'second']}
                  renderInput={({ showTime }) => (
                    <div class="custom-input">
                      <input value={timeValueTo().label} readOnly />
                      <button onClick={showTime}><CalendarIcon /></button>
                    </div>
                  )}
                  onChange={() => { toggleDeviceEvents }}
                  timeAnalogClockCenterDotClass='text-indigo-500'
                  timeAnalogNumberClass='text-slate-700'
                  timeAnalogClockHandClass='text-indigo-500'
                  timeAnalogWrapperClass='text-slate-700'
                  timePickerBottomAreaClass='text-slate-700'
                  timePickerMeridiemBtnClass='text-slate-700'
                  timePickerTopAreaClass='text-slate-700'
                  timePickerWrapperClass='text-slate-700'
                />
              </div>
            </div>

            {/* Events grid */}
            <div class="flex-1 overflow-hidden">
              <div class="ag-theme-alpine h-full w-full">
                {(area() == 'device_events') && (
                  <AgGridSolid
                    rowData={deviceEventsToggle()}
                    columnDefs={columnDefsEvent}
                    defaultColDef={defaultColDef}
                    animateRows={true}
                    autoSizeStrategy={gridOptions.autoSizeStrategy}
                  />
                )}
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default App;
