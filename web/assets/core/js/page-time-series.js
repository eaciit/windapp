'use strict';

viewModel.TurbineHealth = new Object();
var pg = viewModel.TurbineHealth;

vm.currentMenu('Time Series Plots');
vm.currentTitle('Time Series Plots');
vm.breadcrumb([{ title: 'Analysis Tool Box', href: '#' }, { title: 'Time Series Plots', href: viewModel.appName + 'page/timeseries' }]);

pg.chartWindSpeed = function(){
	$("#chartWindSpeed").kendoStockChart({
	  title: {
        text: "Time Series Plots for Wind Speed",
        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: [
        	{
					"date": "2013-12-12",
					"revenue": 1361,
					"costs": 845,
					"income": 516,
					"dashed": 9422
				},
				{
					"date": "2013-12-13",
					"revenue": 4535,
					"costs": 1688,
					"income": 2847,
					"dashed": 7742
				},
				{
					"date": "2013-12-14",
					"revenue": 2466,
					"costs": 117,
					"income": 2349,
					"dashed": 7945
				},
				{
					"date": "2013-12-15",
					"revenue": 200,
					"costs": 1571,
					"income": -1371,
					"dashed": 1262
				},
				{
					"date": "2013-12-16",
					"revenue": 3330,
					"costs": 1552,
					"income": 1778,
					"dashed": 9551
				},
				{
					"date": "2013-12-17",
					"revenue": 3033,
					"costs": 875,
					"income": 2158,
					"dashed": 573
				},
				{
					"date": "2013-12-18",
					"revenue": 3247,
					"costs": 658,
					"income": 2589,
					"dashed": 8417
				},
				{
					"date": "2013-12-19",
					"revenue": 1189,
					"costs": 223,
					"income": 966,
					"dashed": 363
				},
				{
					"date": "2013-12-20",
					"revenue": 544,
					"costs": 391,
					"income": 153,
					"dashed": 1322
				},
				{
					"date": "2013-12-21",
					"revenue": 3540,
					"costs": 160,
					"income": 3380,
					"dashed": 8258
				},
				{
					"date": "2013-12-22",
					"revenue": 67,
					"costs": 992,
					"income": -925,
					"dashed": 8542
				},
				{
					"date": "2013-12-23",
					"revenue": 410,
					"costs": 796,
					"income": -386,
					"dashed": 5322
				},
				{
					"date": "2013-12-24",
					"revenue": 3316,
					"costs": 1821,
					"income": 1495,
					"dashed": 6134
				},
				{
					"date": "2013-12-25",
					"revenue": 1197,
					"costs": 658,
					"income": 539,
					"dashed": 8985
				},
				{
					"date": "2013-12-26",
					"revenue": 4319,
					"costs": 568,
					"income": 3751,
					"dashed": 7767
				},
				{
					"date": "2013-12-27",
					"revenue": 4618,
					"costs": 1781,
					"income": 2837,
					"dashed": 6001
				},
				{
					"date": "2013-12-28",
					"revenue": 2947,
					"costs": 1179,
					"income": 1768,
					"dashed": 692
				},
				{
					"date": "2013-12-29",
					"revenue": 4670,
					"costs": 1974,
					"income": 2696,
					"dashed": 3466
				},
				{
					"date": "2013-12-30",
					"revenue": 2139,
					"costs": 334,
					"income": 1805,
					"dashed": 4114
				},
				{
					"date": "2013-12-31",
					"revenue": 2018,
					"costs": 1175,
					"income": 843,
					"dashed": 5681
				},
				{
					"date": "2014-01-01",
					"revenue": 1682,
					"costs": 338,
					"income": 1344,
					"dashed": 1672
				},
				{
					"date": "2014-01-02",
					"revenue": 2434,
					"costs": 946,
					"income": 1488,
					"dashed": 959
				},
				{
					"date": "2014-01-03",
					"revenue": 60,
					"costs": 50,
					"income": 10,
					"dashed": 837
				},
				{
					"date": "2014-01-04",
					"revenue": 1627,
					"costs": 2248,
					"income": -621,
					"dashed": 892
				},
				{
					"date": "2014-01-05",
					"revenue": 4505,
					"costs": 282,
					"income": 4223,
					"dashed": 3200
				},
				{
					"date": "2014-01-06",
					"revenue": 1215,
					"costs": 2489,
					"income": -1274,
					"dashed": 6691
				},
				{
					"date": "2014-01-07",
					"revenue": 3553,
					"costs": 900,
					"income": 2653,
					"dashed": 6233
				},
				{
					"date": "2014-01-08",
					"revenue": 1942,
					"costs": 57,
					"income": 1885,
					"dashed": 6690
				},
				{
					"date": "2014-01-09",
					"revenue": 1965,
					"costs": 52,
					"income": 1913,
					"dashed": 9926
				},
				{
					"date": "2014-01-10",
					"revenue": 3013,
					"costs": 1039,
					"income": 1974,
					"dashed": 8618
				},
				{
					"date": "2014-01-11",
					"revenue": 629,
					"costs": 21,
					"income": 608,
					"dashed": 9149
				},
				{
					"date": "2014-01-12",
					"revenue": 1419,
					"costs": 741,
					"income": 678,
					"dashed": 8620
				},
				{
					"date": "2014-01-13",
					"revenue": 4321,
					"costs": 1428,
					"income": 2893,
					"dashed": 5730
				},
				{
					"date": "2014-01-14",
					"revenue": 1575,
					"costs": 1396,
					"income": 179,
					"dashed": 4778
				},
				{
					"date": "2014-01-15",
					"revenue": 523,
					"costs": 1482,
					"income": -959,
					"dashed": 5581
				},
				{
					"date": "2014-01-16",
					"revenue": 2093,
					"costs": 730,
					"income": 1363,
					"dashed": 4562
				},
				{
					"date": "2014-01-17",
					"revenue": 1681,
					"costs": 353,
					"income": 1328,
					"dashed": 6335
				},
				{
					"date": "2014-01-18",
					"revenue": 3925,
					"costs": 961,
					"income": 2964,
					"dashed": 1976
				},
				{
					"date": "2014-01-19",
					"revenue": 4458,
					"costs": 1004,
					"income": 3454,
					"dashed": 1610
				},
				{
					"date": "2014-01-20",
					"revenue": 3002,
					"costs": 493,
					"income": 2509,
					"dashed": 1565
				},
				{
					"date": "2014-01-21",
					"revenue": 2176,
					"costs": 659,
					"income": 1517,
					"dashed": 6315
				},
				{
					"date": "2014-01-22",
					"revenue": 2763,
					"costs": 444,
					"income": 2319,
					"dashed": 4995
				},
				{
					"date": "2014-01-23",
					"revenue": 3871,
					"costs": 2006,
					"income": 1865,
					"dashed": 9633
				},
				{
					"date": "2014-01-24",
					"revenue": 2740,
					"costs": 1052,
					"income": 1688,
					"dashed": 9702
				},
				{
					"date": "2014-01-25",
					"revenue": 4735,
					"costs": 2312,
					"income": 2423,
					"dashed": 3952
				},
				{
					"date": "2014-01-26",
					"revenue": 3082,
					"costs": 1229,
					"income": 1853,
					"dashed": 9761
				},
				{
					"date": "2014-01-27",
					"revenue": 2374,
					"costs": 244,
					"income": 2130,
					"dashed": 8020
				},
				{
					"date": "2014-01-28",
					"revenue": 2168,
					"costs": 1358,
					"income": 810,
					"dashed": 1447
				},
				{
					"date": "2014-01-29",
					"revenue": 3363,
					"costs": 976,
					"income": 2387,
					"dashed": 2568
				},
				{
					"date": "2014-01-30",
					"revenue": 4465,
					"costs": 751,
					"income": 3714,
					"dashed": 2608
				},
				{
					"date": "2014-01-31",
					"revenue": 4422,
					"costs": 1139,
					"income": 3283,
					"dashed": 4640
				},
				{
					"date": "2014-02-01",
					"revenue": 797,
					"costs": 101,
					"income": 696,
					"dashed": 6379
				},
				{
					"date": "2014-02-02",
					"revenue": 221,
					"costs": 2284,
					"income": -2063,
					"dashed": 3477
				},
				{
					"date": "2014-02-03",
					"revenue": 412,
					"costs": 1312,
					"income": -900,
					"dashed": 7269
				},
				{
					"date": "2014-02-04",
					"revenue": 4002,
					"costs": 1926,
					"income": 2076,
					"dashed": 1031
				},
				{
					"date": "2014-02-05",
					"revenue": 4407,
					"costs": 1803,
					"income": 2604,
					"dashed": 510
				},
				{
					"date": "2014-02-06",
					"revenue": 4664,
					"costs": 2349,
					"income": 2315,
					"dashed": 2473
				},
				{
					"date": "2014-02-07",
					"revenue": 3057,
					"costs": 2467,
					"income": 590,
					"dashed": 1298
				},
				{
					"date": "2014-02-08",
					"revenue": 865,
					"costs": 1836,
					"income": -971,
					"dashed": 5659
				},
				{
					"date": "2014-02-09",
					"revenue": 1663,
					"costs": 993,
					"income": 670,
					"dashed": 943
				},
				{
					"date": "2014-02-10",
					"revenue": 4458,
					"costs": 1432,
					"income": 3026,
					"dashed": 9764
				},
				{
					"date": "2014-02-11",
					"revenue": 4421,
					"costs": 2235,
					"income": 2186,
					"dashed": 6007
				},
				{
					"date": "2014-02-12",
					"revenue": 4039,
					"costs": 2369,
					"income": 1670,
					"dashed": 3952
				},
				{
					"date": "2014-02-13",
					"revenue": 591,
					"costs": 493,
					"income": 98,
					"dashed": 5562
				},
				{
					"date": "2014-02-14",
					"revenue": 3635,
					"costs": 48,
					"income": 3587,
					"dashed": 5958
				},
				{
					"date": "2014-02-15",
					"revenue": 503,
					"costs": 239,
					"income": 264,
					"dashed": 358
				},
				{
					"date": "2014-02-16",
					"revenue": 2485,
					"costs": 972,
					"income": 1513,
					"dashed": 5499
				},
				{
					"date": "2014-02-17",
					"revenue": 1978,
					"costs": 1075,
					"income": 903,
					"dashed": 2436
				},
				{
					"date": "2014-02-18",
					"revenue": 4399,
					"costs": 2262,
					"income": 2137,
					"dashed": 1856
				},
				{
					"date": "2014-02-19",
					"revenue": 1469,
					"costs": 1523,
					"income": -54,
					"dashed": 6902
				},
				{
					"date": "2014-02-20",
					"revenue": 3361,
					"costs": 330,
					"income": 3031,
					"dashed": 2823
				},
				{
					"date": "2014-02-21",
					"revenue": 2487,
					"costs": 1476,
					"income": 1011,
					"dashed": 8635
				},
				{
					"date": "2014-02-22",
					"revenue": 575,
					"costs": 914,
					"income": -339,
					"dashed": 776
				},
				{
					"date": "2014-02-23",
					"revenue": 3083,
					"costs": 2046,
					"income": 1037,
					"dashed": 9008
				},
				{
					"date": "2014-02-24",
					"revenue": 855,
					"costs": 1412,
					"income": -557,
					"dashed": 374
				},
				{
					"date": "2014-02-25",
					"revenue": 3996,
					"costs": 1527,
					"income": 2469,
					"dashed": 7505
				},
				{
					"date": "2014-02-26",
					"revenue": 3428,
					"costs": 766,
					"income": 2662,
					"dashed": 2680
				},
				{
					"date": "2014-02-27",
					"revenue": 4751,
					"costs": 2131,
					"income": 2620,
					"dashed": 3083
				},
				{
					"date": "2014-02-28",
					"revenue": 1706,
					"costs": 1647,
					"income": 59,
					"dashed": 8174
				},
				{
					"date": "2014-03-01",
					"revenue": 4853,
					"costs": 92,
					"income": 4761,
					"dashed": 7450
				},
				{
					"date": "2014-03-02",
					"revenue": 4856,
					"costs": 767,
					"income": 4089,
					"dashed": 5105
				},
				{
					"date": "2014-03-03",
					"revenue": 4441,
					"costs": 1445,
					"income": 2996,
					"dashed": 6846
				},
				{
					"date": "2014-03-04",
					"revenue": 2990,
					"costs": 591,
					"income": 2399,
					"dashed": 8790
				},
				{
					"date": "2014-03-05",
					"revenue": 2309,
					"costs": 2263,
					"income": 46,
					"dashed": 3573
				},
				{
					"date": "2014-03-06",
					"revenue": 3564,
					"costs": 573,
					"income": 2991,
					"dashed": 8374
				},
				{
					"date": "2014-03-07",
					"revenue": 4623,
					"costs": 1463,
					"income": 3160,
					"dashed": 1079
				},
				{
					"date": "2014-03-08",
					"revenue": 4141,
					"costs": 1321,
					"income": 2820,
					"dashed": 4566
				},
				{
					"date": "2014-03-09",
					"revenue": 1651,
					"costs": 1871,
					"income": -220,
					"dashed": 3839
				},
				{
					"date": "2014-03-10",
					"revenue": 2458,
					"costs": 242,
					"income": 2216,
					"dashed": 2835
				},
				{
					"date": "2014-03-11",
					"revenue": 2550,
					"costs": 731,
					"income": 1819,
					"dashed": 4691
				},
				{
					"date": "2014-03-12",
					"revenue": 2795,
					"costs": 584,
					"income": 2211,
					"dashed": 6175
				},
				{
					"date": "2014-03-13",
					"revenue": 3506,
					"costs": 380,
					"income": 3126,
					"dashed": 7159
				},
				{
					"date": "2014-03-14",
					"revenue": 1699,
					"costs": 428,
					"income": 1271,
					"dashed": 9150
				},
				{
					"date": "2014-03-15",
					"revenue": 188,
					"costs": 1430,
					"income": -1242,
					"dashed": 1296
				},
				{
					"date": "2014-03-16",
					"revenue": 158,
					"costs": 835,
					"income": -677,
					"dashed": 1949
				},
				{
					"date": "2014-03-17",
					"revenue": 516,
					"costs": 2173,
					"income": -1657,
					"dashed": 6999
				},
				{
					"date": "2014-03-18",
					"revenue": 1769,
					"costs": 2074,
					"income": -305,
					"dashed": 2283
				},
				{
					"date": "2014-03-19",
					"revenue": 3042,
					"costs": 1934,
					"income": 1108,
					"dashed": 8095
				},
				{
					"date": "2014-03-20",
					"revenue": 1204,
					"costs": 1166,
					"income": 38,
					"dashed": 8771
				},
				{
					"date": "2014-03-21",
					"revenue": 4330,
					"costs": 20,
					"income": 4310,
					"dashed": 889
				},
				{
					"date": "2014-03-22",
					"revenue": 772,
					"costs": 112,
					"income": 660,
					"dashed": 3131
				},
				{
					"date": "2014-03-23",
					"revenue": 3853,
					"costs": 2140,
					"income": 1713,
					"dashed": 5781
				},
				{
					"date": "2014-03-24",
					"revenue": 3324,
					"costs": 1843,
					"income": 1481,
					"dashed": 8743
				},
				{
					"date": "2014-03-25",
					"revenue": 3165,
					"costs": 1487,
					"income": 1678,
					"dashed": 5554
				},
				{
					"date": "2014-03-26",
					"revenue": 4756,
					"costs": 2091,
					"income": 2665,
					"dashed": 644
				},
				{
					"date": "2014-03-27",
					"revenue": 2311,
					"costs": 437,
					"income": 1874,
					"dashed": 286
				},
				{
					"date": "2014-03-28",
					"revenue": 2511,
					"costs": 1323,
					"income": 1188,
					"dashed": 9027
				},
				{
					"date": "2014-03-29",
					"revenue": 1426,
					"costs": 1144,
					"income": 282,
					"dashed": 7
				},
				{
					"date": "2014-03-30",
					"revenue": 3596,
					"costs": 564,
					"income": 3032,
					"dashed": 4450
				},
				{
					"date": "2014-03-31",
					"revenue": 3953,
					"costs": 1711,
					"income": 2242,
					"dashed": 3910
				},
				{
					"date": "2014-04-01",
					"revenue": 3176,
					"costs": 1525,
					"income": 1651,
					"dashed": 7195
				},
				{
					"date": "2014-04-02",
					"revenue": 4786,
					"costs": 824,
					"income": 3962,
					"dashed": 4465
				},
				{
					"date": "2014-04-03",
					"revenue": 417,
					"costs": 1595,
					"income": -1178,
					"dashed": 3730
				},
				{
					"date": "2014-04-04",
					"revenue": 3202,
					"costs": 2232,
					"income": 970,
					"dashed": 9995
				},
				{
					"date": "2014-04-05",
					"revenue": 833,
					"costs": 2173,
					"income": -1340,
					"dashed": 1780
				},
				{
					"date": "2014-04-06",
					"revenue": 1055,
					"costs": 2310,
					"income": -1255,
					"dashed": 3323
				},
				{
					"date": "2014-04-07",
					"revenue": 2251,
					"costs": 1743,
					"income": 508,
					"dashed": 5952
				},
				{
					"date": "2014-04-08",
					"revenue": 190,
					"costs": 2454,
					"income": -2264,
					"dashed": 2579
				},
				{
					"date": "2014-04-09",
					"revenue": 2615,
					"costs": 1591,
					"income": 1024,
					"dashed": 8140
				},
				{
					"date": "2014-04-10",
					"revenue": 4070,
					"costs": 538,
					"income": 3532,
					"dashed": 9050
				},
				{
					"date": "2014-04-11",
					"revenue": 985,
					"costs": 782,
					"income": 203,
					"dashed": 5413
				},
				{
					"date": "2014-04-12",
					"revenue": 4004,
					"costs": 1071,
					"income": 2933,
					"dashed": 9431
				},
				{
					"date": "2014-04-13",
					"revenue": 4360,
					"costs": 25,
					"income": 4335,
					"dashed": 9327
				},
				{
					"date": "2014-04-14",
					"revenue": 1640,
					"costs": 628,
					"income": 1012,
					"dashed": 1480
				},
				{
					"date": "2014-04-15",
					"revenue": 1742,
					"costs": 1206,
					"income": 536,
					"dashed": 9115
				},
				{
					"date": "2014-04-16",
					"revenue": 4905,
					"costs": 2099,
					"income": 2806,
					"dashed": 4785
				},
				{
					"date": "2014-04-17",
					"revenue": 2569,
					"costs": 1470,
					"income": 1099,
					"dashed": 7490
				},
				{
					"date": "2014-04-18",
					"revenue": 842,
					"costs": 2332,
					"income": -1490,
					"dashed": 3733
				},
				{
					"date": "2014-04-19",
					"revenue": 2332,
					"costs": 533,
					"income": 1799,
					"dashed": 429
				},
				{
					"date": "2014-04-20",
					"revenue": 380,
					"costs": 241,
					"income": 139,
					"dashed": 6214
				},
				{
					"date": "2014-04-21",
					"revenue": 3598,
					"costs": 2213,
					"income": 1385,
					"dashed": 6863
				},
				{
					"date": "2014-04-22",
					"revenue": 2240,
					"costs": 2297,
					"income": -57,
					"dashed": 915
				},
				{
					"date": "2014-04-23",
					"revenue": 4232,
					"costs": 2406,
					"income": 1826,
					"dashed": 1219
				},
				{
					"date": "2014-04-24",
					"revenue": 2615,
					"costs": 2180,
					"income": 435,
					"dashed": 9037
				},
				{
					"date": "2014-04-25",
					"revenue": 3448,
					"costs": 2465,
					"income": 983,
					"dashed": 6075
				},
				{
					"date": "2014-04-26",
					"revenue": 4530,
					"costs": 2193,
					"income": 2337,
					"dashed": 6423
				},
				{
					"date": "2014-04-27",
					"revenue": 4194,
					"costs": 716,
					"income": 3478,
					"dashed": 9983
				},
				{
					"date": "2014-04-28",
					"revenue": 96,
					"costs": 423,
					"income": -327,
					"dashed": 6826
				},
				{
					"date": "2014-04-29",
					"revenue": 4358,
					"costs": 912,
					"income": 3446,
					"dashed": 1479
				},
				{
					"date": "2014-04-30",
					"revenue": 3930,
					"costs": 1992,
					"income": 1938,
					"dashed": 9882
				},
				{
					"date": "2014-05-01",
					"revenue": 4802,
					"costs": 2128,
					"income": 2674,
					"dashed": 5589
				},
				{
					"date": "2014-05-02",
					"revenue": 1850,
					"costs": 862,
					"income": 988,
					"dashed": 1525
				},
				{
					"date": "2014-05-03",
					"revenue": 3110,
					"costs": 1107,
					"income": 2003,
					"dashed": 9486
				},
				{
					"date": "2014-05-04",
					"revenue": 2198,
					"costs": 2082,
					"income": 116,
					"dashed": 851
				},
				{
					"date": "2014-05-05",
					"revenue": 1930,
					"costs": 1865,
					"income": 65,
					"dashed": 4906
				},
				{
					"date": "2014-05-06",
					"revenue": 291,
					"costs": 1253,
					"income": -962,
					"dashed": 7060
				},
				{
					"date": "2014-05-07",
					"revenue": 1935,
					"costs": 31,
					"income": 1904,
					"dashed": 9193
				},
				{
					"date": "2014-05-08",
					"revenue": 2848,
					"costs": 2414,
					"income": 434,
					"dashed": 8282
				},
				{
					"date": "2014-05-09",
					"revenue": 4672,
					"costs": 1807,
					"income": 2865,
					"dashed": 9648
				},
				{
					"date": "2014-05-10",
					"revenue": 1743,
					"costs": 797,
					"income": 946,
					"dashed": 65
				},
				{
					"date": "2014-05-11",
					"revenue": 2532,
					"costs": 940,
					"income": 1592,
					"dashed": 7348
				},
				{
					"date": "2014-05-12",
					"revenue": 4098,
					"costs": 2284,
					"income": 1814,
					"dashed": 6733
				},
				{
					"date": "2014-05-13",
					"revenue": 4567,
					"costs": 2236,
					"income": 2331,
					"dashed": 691
				},
				{
					"date": "2014-05-14",
					"revenue": 109,
					"costs": 2328,
					"income": -2219,
					"dashed": 726
				},
				{
					"date": "2014-05-15",
					"revenue": 4303,
					"costs": 1449,
					"income": 2854,
					"dashed": 5052
				},
				{
					"date": "2014-05-16",
					"revenue": 3745,
					"costs": 1140,
					"income": 2605,
					"dashed": 3005
				},
				{
					"date": "2014-05-17",
					"revenue": 2940,
					"costs": 472,
					"income": 2468,
					"dashed": 9008
				},
				{
					"date": "2014-05-18",
					"revenue": 4208,
					"costs": 1260,
					"income": 2948,
					"dashed": 2308
				},
				{
					"date": "2014-05-19",
					"revenue": 3776,
					"costs": 725,
					"income": 3051,
					"dashed": 8398
				},
				{
					"date": "2014-05-20",
					"revenue": 501,
					"costs": 539,
					"income": -38,
					"dashed": 7182
				},
				{
					"date": "2014-05-21",
					"revenue": 1944,
					"costs": 503,
					"income": 1441,
					"dashed": 7459
				},
				{
					"date": "2014-05-22",
					"revenue": 3216,
					"costs": 482,
					"income": 2734,
					"dashed": 6491
				},
				{
					"date": "2014-05-23",
					"revenue": 3428,
					"costs": 493,
					"income": 2935,
					"dashed": 2900
				},
				{
					"date": "2014-05-24",
					"revenue": 1497,
					"costs": 1716,
					"income": -219,
					"dashed": 7982
				},
				{
					"date": "2014-05-25",
					"revenue": 3329,
					"costs": 266,
					"income": 3063,
					"dashed": 7252
				},
				{
					"date": "2014-05-26",
					"revenue": 2890,
					"costs": 1710,
					"income": 1180,
					"dashed": 8924
				},
				{
					"date": "2014-05-27",
					"revenue": 3381,
					"costs": 2208,
					"income": 1173,
					"dashed": 2368
				},
				{
					"date": "2014-05-28",
					"revenue": 2401,
					"costs": 1556,
					"income": 845,
					"dashed": 4067
				},
				{
					"date": "2014-05-29",
					"revenue": 1423,
					"costs": 726,
					"income": 697,
					"dashed": 9387
				},
				{
					"date": "2014-05-30",
					"revenue": 2523,
					"costs": 171,
					"income": 2352,
					"dashed": 332
				},
				{
					"date": "2014-05-31",
					"revenue": 3483,
					"costs": 863,
					"income": 2620,
					"dashed": 7330
				},
				{
					"date": "2014-06-01",
					"revenue": 2575,
					"costs": 2113,
					"income": 462,
					"dashed": 9069
				},
				{
					"date": "2014-06-02",
					"revenue": 604,
					"costs": 1802,
					"income": -1198,
					"dashed": 6452
				},
				{
					"date": "2014-06-03",
					"revenue": 2149,
					"costs": 1851,
					"income": 298,
					"dashed": 3410
				},
				{
					"date": "2014-06-04",
					"revenue": 4335,
					"costs": 1376,
					"income": 2959,
					"dashed": 9354
				},
				{
					"date": "2014-06-05",
					"revenue": 3123,
					"costs": 1069,
					"income": 2054,
					"dashed": 3957
				},
				{
					"date": "2014-06-06",
					"revenue": 2475,
					"costs": 2462,
					"income": 13,
					"dashed": 4780
				},
				{
					"date": "2014-06-07",
					"revenue": 2418,
					"costs": 317,
					"income": 2101,
					"dashed": 6667
				},
				{
					"date": "2014-06-08",
					"revenue": 609,
					"costs": 2015,
					"income": -1406,
					"dashed": 1315
				},
				{
					"date": "2014-06-09",
					"revenue": 1246,
					"costs": 1042,
					"income": 204,
					"dashed": 2856
				},
				{
					"date": "2014-06-10",
					"revenue": 571,
					"costs": 641,
					"income": -70,
					"dashed": 7961
				},
				{
					"date": "2014-06-11",
					"revenue": 878,
					"costs": 206,
					"income": 672,
					"dashed": 5588
				},
				{
					"date": "2014-06-12",
					"revenue": 1740,
					"costs": 1129,
					"income": 611,
					"dashed": 8059
				},
				{
					"date": "2014-06-13",
					"revenue": 4143,
					"costs": 2404,
					"income": 1739,
					"dashed": 8723
				},
				{
					"date": "2014-06-14",
					"revenue": 1549,
					"costs": 1758,
					"income": -209,
					"dashed": 7129
				},
				{
					"date": "2014-06-15",
					"revenue": 1690,
					"costs": 2078,
					"income": -388,
					"dashed": 9956
				},
				{
					"date": "2014-06-16",
					"revenue": 3995,
					"costs": 2145,
					"income": 1850,
					"dashed": 1051
				},
				{
					"date": "2014-06-17",
					"revenue": 4382,
					"costs": 672,
					"income": 3710,
					"dashed": 7998
				},
				{
					"date": "2014-06-18",
					"revenue": 32,
					"costs": 560,
					"income": -528,
					"dashed": 7759
				},
				{
					"date": "2014-06-19",
					"revenue": 4593,
					"costs": 1397,
					"income": 3196,
					"dashed": 4978
				},
				{
					"date": "2014-06-20",
					"revenue": 13,
					"costs": 2479,
					"income": -2466,
					"dashed": 4247
				},
				{
					"date": "2014-06-21",
					"revenue": 3080,
					"costs": 584,
					"income": 2496,
					"dashed": 5
				},
				{
					"date": "2014-06-22",
					"revenue": 2174,
					"costs": 2180,
					"income": -6,
					"dashed": 4214
				},
				{
					"date": "2014-06-23",
					"revenue": 572,
					"costs": 362,
					"income": 210,
					"dashed": 460
				},
				{
					"date": "2014-06-24",
					"revenue": 3191,
					"costs": 2081,
					"income": 1110,
					"dashed": 3813
				},
				{
					"date": "2014-06-25",
					"revenue": 1570,
					"costs": 1419,
					"income": 151,
					"dashed": 2496
				},
				{
					"date": "2014-06-26",
					"revenue": 2104,
					"costs": 323,
					"income": 1781,
					"dashed": 5596
				},
				{
					"date": "2014-06-27",
					"revenue": 4694,
					"costs": 717,
					"income": 3977,
					"dashed": 4784
				},
				{
					"date": "2014-06-28",
					"revenue": 1719,
					"costs": 2010,
					"income": -291,
					"dashed": 2923
				},
				{
					"date": "2014-06-29",
					"revenue": 2041,
					"costs": 1391,
					"income": 650,
					"dashed": 7431
				},
				{
					"date": "2014-06-30",
					"revenue": 527,
					"costs": 1307,
					"income": -780,
					"dashed": 4074
				},
				{
					"date": "2014-07-01",
					"revenue": 1916,
					"costs": 2486,
					"income": -570,
					"dashed": 9763
				},
				{
					"date": "2014-07-02",
					"revenue": 4336,
					"costs": 1220,
					"income": 3116,
					"dashed": 5549
				},
				{
					"date": "2014-07-03",
					"revenue": 4208,
					"costs": 1276,
					"income": 2932,
					"dashed": 6095
				},
				{
					"date": "2014-07-04",
					"revenue": 4206,
					"costs": 2458,
					"income": 1748,
					"dashed": 4026
				},
				{
					"date": "2014-07-05",
					"revenue": 3974,
					"costs": 1738,
					"income": 2236,
					"dashed": 7334
				},
				{
					"date": "2014-07-06",
					"revenue": 1112,
					"costs": 295,
					"income": 817,
					"dashed": 7225
				},
				{
					"date": "2014-07-07",
					"revenue": 2279,
					"costs": 345,
					"income": 1934,
					"dashed": 9051
				},
				{
					"date": "2014-07-08",
					"revenue": 916,
					"costs": 2344,
					"income": -1428,
					"dashed": 3660
				},
				{
					"date": "2014-07-09",
					"revenue": 3770,
					"costs": 1591,
					"income": 2179,
					"dashed": 1924
				},
				{
					"date": "2014-07-10",
					"revenue": 1635,
					"costs": 783,
					"income": 852,
					"dashed": 2477
				},
				{
					"date": "2014-07-11",
					"revenue": 4923,
					"costs": 672,
					"income": 4251,
					"dashed": 9553
				},
				{
					"date": "2014-07-12",
					"revenue": 1963,
					"costs": 1833,
					"income": 130,
					"dashed": 9919
				},
				{
					"date": "2014-07-13",
					"revenue": 3543,
					"costs": 2260,
					"income": 1283,
					"dashed": 8570
				},
				{
					"date": "2014-07-14",
					"revenue": 3677,
					"costs": 313,
					"income": 3364,
					"dashed": 3712
				},
				{
					"date": "2014-07-15",
					"revenue": 1650,
					"costs": 176,
					"income": 1474,
					"dashed": 322
				},
				{
					"date": "2014-07-16",
					"revenue": 4760,
					"costs": 314,
					"income": 4446,
					"dashed": 8118
				},
				{
					"date": "2014-07-17",
					"revenue": 2194,
					"costs": 2132,
					"income": 62,
					"dashed": 554
				},
				{
					"date": "2014-07-18",
					"revenue": 4629,
					"costs": 544,
					"income": 4085,
					"dashed": 9034
				},
				{
					"date": "2014-07-19",
					"revenue": 4702,
					"costs": 439,
					"income": 4263,
					"dashed": 7172
				},
				{
					"date": "2014-07-20",
					"revenue": 4816,
					"costs": 138,
					"income": 4678,
					"dashed": 3242
				},
				{
					"date": "2014-07-21",
					"revenue": 1259,
					"costs": 2478,
					"income": -1219,
					"dashed": 7665
				},
				{
					"date": "2014-07-22",
					"revenue": 1947,
					"costs": 91,
					"income": 1856,
					"dashed": 8381
				},
				{
					"date": "2014-07-23",
					"revenue": 2921,
					"costs": 828,
					"income": 2093,
					"dashed": 7873
				},
				{
					"date": "2014-07-24",
					"revenue": 4958,
					"costs": 2462,
					"income": 2496,
					"dashed": 6109
				},
				{
					"date": "2014-07-25",
					"revenue": 1720,
					"costs": 1623,
					"income": 97,
					"dashed": 8819
				},
				{
					"date": "2014-07-26",
					"revenue": 4272,
					"costs": 437,
					"income": 3835,
					"dashed": 4096
				},
				{
					"date": "2014-07-27",
					"revenue": 3641,
					"costs": 2344,
					"income": 1297,
					"dashed": 9577
				},
				{
					"date": "2014-07-28",
					"revenue": 4571,
					"costs": 654,
					"income": 3917,
					"dashed": 8119
				},
				{
					"date": "2014-07-29",
					"revenue": 4714,
					"costs": 2251,
					"income": 2463,
					"dashed": 4181
				},
				{
					"date": "2014-07-30",
					"revenue": 1731,
					"costs": 802,
					"income": 929,
					"dashed": 3577
				},
				{
					"date": "2014-07-31",
					"revenue": 3682,
					"costs": 128,
					"income": 3554,
					"dashed": 2910
				},
				{
					"date": "2014-08-01",
					"revenue": 259,
					"costs": 1574,
					"income": -1315,
					"dashed": 4107
				},
				{
					"date": "2014-08-02",
					"revenue": 4094,
					"costs": 2444,
					"income": 1650,
					"dashed": 2883
				},
				{
					"date": "2014-08-03",
					"revenue": 1420,
					"costs": 1724,
					"income": -304,
					"dashed": 7104
				},
				{
					"date": "2014-08-04",
					"revenue": 4939,
					"costs": 2038,
					"income": 2901,
					"dashed": 9687
				},
				{
					"date": "2014-08-05",
					"revenue": 3809,
					"costs": 1969,
					"income": 1840,
					"dashed": 5870
				},
				{
					"date": "2014-08-06",
					"revenue": 3512,
					"costs": 1451,
					"income": 2061,
					"dashed": 4681
				},
				{
					"date": "2014-08-07",
					"revenue": 3695,
					"costs": 837,
					"income": 2858,
					"dashed": 4969
				},
				{
					"date": "2014-08-08",
					"revenue": 2638,
					"costs": 1023,
					"income": 1615,
					"dashed": 133
				},
				{
					"date": "2014-08-09",
					"revenue": 4425,
					"costs": 448,
					"income": 3977,
					"dashed": 2727
				},
				{
					"date": "2014-08-10",
					"revenue": 2805,
					"costs": 1862,
					"income": 943,
					"dashed": 3443
				},
				{
					"date": "2014-08-11",
					"revenue": 427,
					"costs": 2172,
					"income": -1745,
					"dashed": 8399
				},
				{
					"date": "2014-08-12",
					"revenue": 3791,
					"costs": 1787,
					"income": 2004,
					"dashed": 3421
				},
				{
					"date": "2014-08-13",
					"revenue": 1453,
					"costs": 1617,
					"income": -164,
					"dashed": 1017
				},
				{
					"date": "2014-08-14",
					"revenue": 671,
					"costs": 2321,
					"income": -1650,
					"dashed": 3321
				},
				{
					"date": "2014-08-15",
					"revenue": 2853,
					"costs": 1221,
					"income": 1632,
					"dashed": 263
				},
				{
					"date": "2014-08-16",
					"revenue": 4705,
					"costs": 1339,
					"income": 3366,
					"dashed": 8853
				},
				{
					"date": "2014-08-17",
					"revenue": 2182,
					"costs": 64,
					"income": 2118,
					"dashed": 1777
				},
				{
					"date": "2014-08-18",
					"revenue": 715,
					"costs": 1510,
					"income": -795,
					"dashed": 2666
				},
				{
					"date": "2014-08-19",
					"revenue": 389,
					"costs": 2235,
					"income": -1846,
					"dashed": 8791
				},
				{
					"date": "2014-08-20",
					"revenue": 490,
					"costs": 628,
					"income": -138,
					"dashed": 3944
				},
				{
					"date": "2014-08-21",
					"revenue": 4438,
					"costs": 795,
					"income": 3643,
					"dashed": 9156
				},
				{
					"date": "2014-08-22",
					"revenue": 320,
					"costs": 535,
					"income": -215,
					"dashed": 5232
				},
				{
					"date": "2014-08-23",
					"revenue": 177,
					"costs": 1898,
					"income": -1721,
					"dashed": 3998
				},
				{
					"date": "2014-08-24",
					"revenue": 3184,
					"costs": 1016,
					"income": 2168,
					"dashed": 322
				},
				{
					"date": "2014-08-25",
					"revenue": 577,
					"costs": 350,
					"income": 227,
					"dashed": 6485
				},
				{
					"date": "2014-08-26",
					"revenue": 277,
					"costs": 1013,
					"income": -736,
					"dashed": 3394
				},
				{
					"date": "2014-08-27",
					"revenue": 414,
					"costs": 1136,
					"income": -722,
					"dashed": 9021
				},
				{
					"date": "2014-08-28",
					"revenue": 1505,
					"costs": 1285,
					"income": 220,
					"dashed": 4713
				},
				{
					"date": "2014-08-29",
					"revenue": 921,
					"costs": 816,
					"income": 105,
					"dashed": 6951
				},
				{
					"date": "2014-08-30",
					"revenue": 646,
					"costs": 574,
					"income": 72,
					"dashed": 2287
				},
				{
					"date": "2014-08-31",
					"revenue": 3473,
					"costs": 1899,
					"income": 1574,
					"dashed": 3700
				},
				{
					"date": "2014-09-01",
					"revenue": 3082,
					"costs": 2074,
					"income": 1008,
					"dashed": 631
				},
				{
					"date": "2014-09-02",
					"revenue": 3200,
					"costs": 120,
					"income": 3080,
					"dashed": 7341
				},
				{
					"date": "2014-09-03",
					"revenue": 3822,
					"costs": 755,
					"income": 3067,
					"dashed": 5715
				},
				{
					"date": "2014-09-04",
					"revenue": 4717,
					"costs": 886,
					"income": 3831,
					"dashed": 6248
				},
				{
					"date": "2014-09-05",
					"revenue": 3925,
					"costs": 2178,
					"income": 1747,
					"dashed": 5740
				},
				{
					"date": "2014-09-06",
					"revenue": 3875,
					"costs": 642,
					"income": 3233,
					"dashed": 2951
				},
				{
					"date": "2014-09-07",
					"revenue": 3743,
					"costs": 2076,
					"income": 1667,
					"dashed": 507
				},
				{
					"date": "2014-09-08",
					"revenue": 500,
					"costs": 888,
					"income": -388,
					"dashed": 991
				},
				{
					"date": "2014-09-09",
					"revenue": 4644,
					"costs": 1727,
					"income": 2917,
					"dashed": 5579
				},
				{
					"date": "2014-09-10",
					"revenue": 3096,
					"costs": 382,
					"income": 2714,
					"dashed": 217
				},
				{
					"date": "2014-09-11",
					"revenue": 3683,
					"costs": 232,
					"income": 3451,
					"dashed": 4700
				},
				{
					"date": "2014-09-12",
					"revenue": 1640,
					"costs": 291,
					"income": 1349,
					"dashed": 1498
				},
				{
					"date": "2014-09-13",
					"revenue": 4763,
					"costs": 1117,
					"income": 3646,
					"dashed": 3025
				},
				{
					"date": "2014-09-14",
					"revenue": 1575,
					"costs": 371,
					"income": 1204,
					"dashed": 1385
				},
				{
					"date": "2014-09-15",
					"revenue": 2631,
					"costs": 1568,
					"income": 1063,
					"dashed": 9678
				},
				{
					"date": "2014-09-16",
					"revenue": 4281,
					"costs": 2027,
					"income": 2254,
					"dashed": 7144
				},
				{
					"date": "2014-09-17",
					"revenue": 4647,
					"costs": 1173,
					"income": 3474,
					"dashed": 2098
				},
				{
					"date": "2014-09-18",
					"revenue": 1822,
					"costs": 1791,
					"income": 31,
					"dashed": 4297
				},
				{
					"date": "2014-09-19",
					"revenue": 761,
					"costs": 1487,
					"income": -726,
					"dashed": 9044
				},
				{
					"date": "2014-09-20",
					"revenue": 2685,
					"costs": 729,
					"income": 1956,
					"dashed": 616
				},
				{
					"date": "2014-09-21",
					"revenue": 1661,
					"costs": 773,
					"income": 888,
					"dashed": 3349
				},
				{
					"date": "2014-09-22",
					"revenue": 2743,
					"costs": 17,
					"income": 2726,
					"dashed": 8259
				},
				{
					"date": "2014-09-23",
					"revenue": 4608,
					"costs": 1345,
					"income": 3263,
					"dashed": 1042
				},
				{
					"date": "2014-09-24",
					"revenue": 1267,
					"costs": 1301,
					"income": -34,
					"dashed": 3088
				},
				{
					"date": "2014-09-25",
					"revenue": 2127,
					"costs": 1726,
					"income": 401,
					"dashed": 2334
				},
				{
					"date": "2014-09-26",
					"revenue": 2271,
					"costs": 2281,
					"income": -10,
					"dashed": 5247
				},
				{
					"date": "2014-09-27",
					"revenue": 1383,
					"costs": 1203,
					"income": 180,
					"dashed": 1510
				},
				{
					"date": "2014-09-28",
					"revenue": 2965,
					"costs": 1175,
					"income": 1790,
					"dashed": 4496
				},
				{
					"date": "2014-09-29",
					"revenue": 3334,
					"costs": 271,
					"income": 3063,
					"dashed": 2382
				},
				{
					"date": "2014-09-30",
					"revenue": 278,
					"costs": 421,
					"income": -143,
					"dashed": 3750
				},
				{
					"date": "2014-10-01",
					"revenue": 1989,
					"costs": 1165,
					"income": 824,
					"dashed": 7319
				},
				{
					"date": "2014-10-02",
					"revenue": 3273,
					"costs": 1817,
					"income": 1456,
					"dashed": 3268
				},
				{
					"date": "2014-10-03",
					"revenue": 4728,
					"costs": 1219,
					"income": 3509,
					"dashed": 6640
				},
				{
					"date": "2014-10-04",
					"revenue": 2760,
					"costs": 1241,
					"income": 1519,
					"dashed": 2925
				},
				{
					"date": "2014-10-05",
					"revenue": 4124,
					"costs": 977,
					"income": 3147,
					"dashed": 4482
				},
				{
					"date": "2014-10-06",
					"revenue": 1752,
					"costs": 1044,
					"income": 708,
					"dashed": 4413
				},
				{
					"date": "2014-10-07",
					"revenue": 205,
					"costs": 941,
					"income": -736,
					"dashed": 3880
				},
				{
					"date": "2014-10-08",
					"revenue": 4007,
					"costs": 337,
					"income": 3670,
					"dashed": 7805
				},
				{
					"date": "2014-10-09",
					"revenue": 3960,
					"costs": 1273,
					"income": 2687,
					"dashed": 8150
				},
				{
					"date": "2014-10-10",
					"revenue": 3874,
					"costs": 1121,
					"income": 2753,
					"dashed": 543
				},
				{
					"date": "2014-10-11",
					"revenue": 4501,
					"costs": 361,
					"income": 4140,
					"dashed": 1586
				},
				{
					"date": "2014-10-12",
					"revenue": 159,
					"costs": 54,
					"income": 105,
					"dashed": 8787
				},
				{
					"date": "2014-10-13",
					"revenue": 3774,
					"costs": 1413,
					"income": 2361,
					"dashed": 1728
				},
				{
					"date": "2014-10-14",
					"revenue": 1197,
					"costs": 1108,
					"income": 89,
					"dashed": 279
				},
				{
					"date": "2014-10-15",
					"revenue": 530,
					"costs": 1703,
					"income": -1173,
					"dashed": 9112
				},
				{
					"date": "2014-10-16",
					"revenue": 1886,
					"costs": 1791,
					"income": 95,
					"dashed": 7297
				},
				{
					"date": "2014-10-17",
					"revenue": 1837,
					"costs": 689,
					"income": 1148,
					"dashed": 45
				},
				{
					"date": "2014-10-18",
					"revenue": 3739,
					"costs": 1873,
					"income": 1866,
					"dashed": 1318
				},
				{
					"date": "2014-10-19",
					"revenue": 4353,
					"costs": 959,
					"income": 3394,
					"dashed": 2017
				},
				{
					"date": "2014-10-20",
					"revenue": 1516,
					"costs": 1917,
					"income": -401,
					"dashed": 780
				},
				{
					"date": "2014-10-21",
					"revenue": 3461,
					"costs": 1665,
					"income": 1796,
					"dashed": 4965
				},
				{
					"date": "2014-10-22",
					"revenue": 4079,
					"costs": 1411,
					"income": 2668,
					"dashed": 3458
				},
				{
					"date": "2014-10-23",
					"revenue": 3629,
					"costs": 1727,
					"income": 1902,
					"dashed": 81
				},
				{
					"date": "2014-10-24",
					"revenue": 3094,
					"costs": 520,
					"income": 2574,
					"dashed": 9382
				},
				{
					"date": "2014-10-25",
					"revenue": 3080,
					"costs": 1953,
					"income": 1127,
					"dashed": 618
				},
				{
					"date": "2014-10-26",
					"revenue": 2693,
					"costs": 1635,
					"income": 1058,
					"dashed": 9149
				},
				{
					"date": "2014-10-27",
					"revenue": 4376,
					"costs": 128,
					"income": 4248,
					"dashed": 3518
				},
				{
					"date": "2014-10-28",
					"revenue": 4135,
					"costs": 1465,
					"income": 2670,
					"dashed": 1582
				},
				{
					"date": "2014-10-29",
					"revenue": 3110,
					"costs": 333,
					"income": 2777,
					"dashed": 3753
				},
				{
					"date": "2014-10-30",
					"revenue": 1185,
					"costs": 2111,
					"income": -926,
					"dashed": 6724
				},
				{
					"date": "2014-10-31",
					"revenue": 693,
					"costs": 1947,
					"income": -1254,
					"dashed": 1931
				},
				{
					"date": "2014-11-01",
					"revenue": 1420,
					"costs": 1240,
					"income": 180,
					"dashed": 5125
				},
				{
					"date": "2014-11-02",
					"revenue": 579,
					"costs": 249,
					"income": 330,
					"dashed": 2232
				},
				{
					"date": "2014-11-03",
					"revenue": 1062,
					"costs": 2396,
					"income": -1334,
					"dashed": 3393
				},
				{
					"date": "2014-11-04",
					"revenue": 3259,
					"costs": 906,
					"income": 2353,
					"dashed": 2581
				},
				{
					"date": "2014-11-05",
					"revenue": 1241,
					"costs": 1560,
					"income": -319,
					"dashed": 164
				},
				{
					"date": "2014-11-06",
					"revenue": 975,
					"costs": 950,
					"income": 25,
					"dashed": 1583
				},
				{
					"date": "2014-11-07",
					"revenue": 2537,
					"costs": 78,
					"income": 2459,
					"dashed": 3674
				},
				{
					"date": "2014-11-08",
					"revenue": 350,
					"costs": 2447,
					"income": -2097,
					"dashed": 5946
				},
				{
					"date": "2014-11-09",
					"revenue": 2692,
					"costs": 215,
					"income": 2477,
					"dashed": 953
				},
				{
					"date": "2014-11-10",
					"revenue": 4231,
					"costs": 1665,
					"income": 2566,
					"dashed": 7036
				},
				{
					"date": "2014-11-11",
					"revenue": 4149,
					"costs": 2074,
					"income": 2075,
					"dashed": 7015
				},
				{
					"date": "2014-11-12",
					"revenue": 3228,
					"costs": 2229,
					"income": 999,
					"dashed": 4754
				},
				{
					"date": "2014-11-13",
					"revenue": 3067,
					"costs": 92,
					"income": 2975,
					"dashed": 6024
				},
				{
					"date": "2014-11-14",
					"revenue": 96,
					"costs": 559,
					"income": -463,
					"dashed": 7169
				},
				{
					"date": "2014-11-15",
					"revenue": 1535,
					"costs": 1762,
					"income": -227,
					"dashed": 2621
				},
				{
					"date": "2014-11-16",
					"revenue": 1968,
					"costs": 2159,
					"income": -191,
					"dashed": 4929
				},
				{
					"date": "2014-11-17",
					"revenue": 5,
					"costs": 167,
					"income": -162,
					"dashed": 7554
				},
				{
					"date": "2014-11-18",
					"revenue": 1770,
					"costs": 1101,
					"income": 669,
					"dashed": 2206
				},
				{
					"date": "2014-11-19",
					"revenue": 1121,
					"costs": 913,
					"income": 208,
					"dashed": 4821
				},
				{
					"date": "2014-11-20",
					"revenue": 1741,
					"costs": 732,
					"income": 1009,
					"dashed": 3593
				},
				{
					"date": "2014-11-21",
					"revenue": 3182,
					"costs": 921,
					"income": 2261,
					"dashed": 8375
				},
				{
					"date": "2014-11-22",
					"revenue": 3463,
					"costs": 342,
					"income": 3121,
					"dashed": 3957
				},
				{
					"date": "2014-11-23",
					"revenue": 3201,
					"costs": 881,
					"income": 2320,
					"dashed": 3537
				},
				{
					"date": "2014-11-24",
					"revenue": 3333,
					"costs": 932,
					"income": 2401,
					"dashed": 9479
				},
				{
					"date": "2014-11-25",
					"revenue": 835,
					"costs": 513,
					"income": 322,
					"dashed": 9799
				},
				{
					"date": "2014-11-26",
					"revenue": 1710,
					"costs": 1811,
					"income": -101,
					"dashed": 1560
				},
				{
					"date": "2014-11-27",
					"revenue": 287,
					"costs": 437,
					"income": -150,
					"dashed": 96
				},
				{
					"date": "2014-11-28",
					"revenue": 2499,
					"costs": 2402,
					"income": 97,
					"dashed": 5028
				},
				{
					"date": "2014-11-29",
					"revenue": 3833,
					"costs": 1766,
					"income": 2067,
					"dashed": 1002
				},
				{
					"date": "2014-11-30",
					"revenue": 1627,
					"costs": 1170,
					"income": 457,
					"dashed": 3079
				},
				{
					"date": "2014-12-01",
					"revenue": 4092,
					"costs": 2155,
					"income": 1937,
					"dashed": 298
				},
				{
					"date": "2014-12-02",
					"revenue": 1372,
					"costs": 463,
					"income": 909,
					"dashed": 5615
				},
				{
					"date": "2014-12-03",
					"revenue": 1479,
					"costs": 539,
					"income": 940,
					"dashed": 9070
				},
				{
					"date": "2014-12-04",
					"revenue": 4886,
					"costs": 1,
					"income": 4885,
					"dashed": 1585
				},
				{
					"date": "2014-12-05",
					"revenue": 4198,
					"costs": 45,
					"income": 4153,
					"dashed": 9867
				},
				{
					"date": "2014-12-06",
					"revenue": 4810,
					"costs": 21,
					"income": 4789,
					"dashed": 2183
				},
				{
					"date": "2014-12-07",
					"revenue": 147,
					"costs": 378,
					"income": -231,
					"dashed": 6553
				},
				{
					"date": "2014-12-08",
					"revenue": 4807,
					"costs": 324,
					"income": 4483,
					"dashed": 8464
				},
				{
					"date": "2014-12-09",
					"revenue": 3359,
					"costs": 1424,
					"income": 1935,
					"dashed": 3404
				},
				{
					"date": "2014-12-10",
					"revenue": 2110,
					"costs": 318,
					"income": 1792,
					"dashed": 7860
				},
				{
					"date": "2014-12-11",
					"revenue": 992,
					"costs": 147,
					"income": 845,
					"dashed": 6405
				},
				{
					"date": "2014-12-12",
					"revenue": 4239,
					"costs": 996,
					"income": 3243,
					"dashed": 2167
				},
				{
					"date": "2014-12-13",
					"revenue": 1025,
					"costs": 1915,
					"income": -890,
					"dashed": 5378
				},
				{
					"date": "2014-12-14",
					"revenue": 3124,
					"costs": 863,
					"income": 2261,
					"dashed": 4270
				},
				{
					"date": "2014-12-15",
					"revenue": 3577,
					"costs": 1448,
					"income": 2129,
					"dashed": 8032
				},
				{
					"date": "2014-12-16",
					"revenue": 3913,
					"costs": 915,
					"income": 2998,
					"dashed": 2181
				},
				{
					"date": "2014-12-17",
					"revenue": 4350,
					"costs": 856,
					"income": 3494,
					"dashed": 7719
				},
				{
					"date": "2014-12-18",
					"revenue": 4528,
					"costs": 2424,
					"income": 2104,
					"dashed": 1834
				},
				{
					"date": "2014-12-19",
					"revenue": 3252,
					"costs": 375,
					"income": 2877,
					"dashed": 5577
				},
				{
					"date": "2014-12-20",
					"revenue": 153,
					"costs": 1327,
					"income": -1174,
					"dashed": 6171
				},
				{
					"date": "2014-12-21",
					"revenue": 1193,
					"costs": 618,
					"income": 575,
					"dashed": 7588
				},
				{
					"date": "2014-12-22",
					"revenue": 726,
					"costs": 961,
					"income": -235,
					"dashed": 2806
				},
				{
					"date": "2014-12-23",
					"revenue": 2245,
					"costs": 2342,
					"income": -97,
					"dashed": 2972
				},
				{
					"date": "2014-12-24",
					"revenue": 4389,
					"costs": 170,
					"income": 4219,
					"dashed": 9920
				},
				{
					"date": "2014-12-25",
					"revenue": 4306,
					"costs": 1806,
					"income": 2500,
					"dashed": 7772
				},
				{
					"date": "2014-12-26",
					"revenue": 4276,
					"costs": 2269,
					"income": 2007,
					"dashed": 3633
				},
				{
					"date": "2014-12-27",
					"revenue": 1136,
					"costs": 1270,
					"income": -134,
					"dashed": 717
				},
				{
					"date": "2014-12-28",
					"revenue": 2750,
					"costs": 1129,
					"income": 1621,
					"dashed": 4484
				},
				{
					"date": "2014-12-29",
					"revenue": 485,
					"costs": 1914,
					"income": -1429,
					"dashed": 1296
				},
				{
					"date": "2014-12-30",
					"revenue": 4992,
					"costs": 1582,
					"income": 3410,
					"dashed": 5231
				}]
      },
      theme: "flat",
      seriesDefaults: {
	        area: {
	            line: {
	                style: "smooth"
	            }
	        }
	    },
      dateField: "date",
      series: [{
        type: "area",
        field: "revenue",
        aggregate: "sum", 
        color: "#337ab7",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        series: [{
          type: "area",
          field: "revenue",
          aggregate: "sum",
          color: "#337ab7",
        }]
      },
      valueAxis: {
        title: {
            text: "Wind Speed",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
	  categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
    });
} 
pg.chartProduction = function(){
	$("#chartProduction").kendoStockChart({
	  title: {
        text: "Time Series Plots for Production",
        font: '14px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif',
      },
      legend: {
        position: "top",
        visible: false
      },
      dataSource: {
        data: [
        	{
					"date": "2013-12-12",
					"revenue": 1361,
					"costs": 845,
					"income": 516,
					"dashed": 9422
				},
				{
					"date": "2013-12-13",
					"revenue": 4535,
					"costs": 1688,
					"income": 2847,
					"dashed": 7742
				},
				{
					"date": "2013-12-14",
					"revenue": 2466,
					"costs": 117,
					"income": 2349,
					"dashed": 7945
				},
				{
					"date": "2013-12-15",
					"revenue": 200,
					"costs": 1571,
					"income": -1371,
					"dashed": 1262
				},
				{
					"date": "2013-12-16",
					"revenue": 3330,
					"costs": 1552,
					"income": 1778,
					"dashed": 9551
				},
				{
					"date": "2013-12-17",
					"revenue": 3033,
					"costs": 875,
					"income": 2158,
					"dashed": 573
				},
				{
					"date": "2013-12-18",
					"revenue": 3247,
					"costs": 658,
					"income": 2589,
					"dashed": 8417
				},
				{
					"date": "2013-12-19",
					"revenue": 1189,
					"costs": 223,
					"income": 966,
					"dashed": 363
				},
				{
					"date": "2013-12-20",
					"revenue": 544,
					"costs": 391,
					"income": 153,
					"dashed": 1322
				},
				{
					"date": "2013-12-21",
					"revenue": 3540,
					"costs": 160,
					"income": 3380,
					"dashed": 8258
				},
				{
					"date": "2013-12-22",
					"revenue": 67,
					"costs": 992,
					"income": -925,
					"dashed": 8542
				},
				{
					"date": "2013-12-23",
					"revenue": 410,
					"costs": 796,
					"income": -386,
					"dashed": 5322
				},
				{
					"date": "2013-12-24",
					"revenue": 3316,
					"costs": 1821,
					"income": 1495,
					"dashed": 6134
				},
				{
					"date": "2013-12-25",
					"revenue": 1197,
					"costs": 658,
					"income": 539,
					"dashed": 8985
				},
				{
					"date": "2013-12-26",
					"revenue": 4319,
					"costs": 568,
					"income": 3751,
					"dashed": 7767
				},
				{
					"date": "2013-12-27",
					"revenue": 4618,
					"costs": 1781,
					"income": 2837,
					"dashed": 6001
				},
				{
					"date": "2013-12-28",
					"revenue": 2947,
					"costs": 1179,
					"income": 1768,
					"dashed": 692
				},
				{
					"date": "2013-12-29",
					"revenue": 4670,
					"costs": 1974,
					"income": 2696,
					"dashed": 3466
				},
				{
					"date": "2013-12-30",
					"revenue": 2139,
					"costs": 334,
					"income": 1805,
					"dashed": 4114
				},
				{
					"date": "2013-12-31",
					"revenue": 2018,
					"costs": 1175,
					"income": 843,
					"dashed": 5681
				},
				{
					"date": "2014-01-01",
					"revenue": 1682,
					"costs": 338,
					"income": 1344,
					"dashed": 1672
				},
				{
					"date": "2014-01-02",
					"revenue": 2434,
					"costs": 946,
					"income": 1488,
					"dashed": 959
				},
				{
					"date": "2014-01-03",
					"revenue": 60,
					"costs": 50,
					"income": 10,
					"dashed": 837
				},
				{
					"date": "2014-01-04",
					"revenue": 1627,
					"costs": 2248,
					"income": -621,
					"dashed": 892
				},
				{
					"date": "2014-01-05",
					"revenue": 4505,
					"costs": 282,
					"income": 4223,
					"dashed": 3200
				},
				{
					"date": "2014-01-06",
					"revenue": 1215,
					"costs": 2489,
					"income": -1274,
					"dashed": 6691
				},
				{
					"date": "2014-01-07",
					"revenue": 3553,
					"costs": 900,
					"income": 2653,
					"dashed": 6233
				},
				{
					"date": "2014-01-08",
					"revenue": 1942,
					"costs": 57,
					"income": 1885,
					"dashed": 6690
				},
				{
					"date": "2014-01-09",
					"revenue": 1965,
					"costs": 52,
					"income": 1913,
					"dashed": 9926
				},
				{
					"date": "2014-01-10",
					"revenue": 3013,
					"costs": 1039,
					"income": 1974,
					"dashed": 8618
				},
				{
					"date": "2014-01-11",
					"revenue": 629,
					"costs": 21,
					"income": 608,
					"dashed": 9149
				},
				{
					"date": "2014-01-12",
					"revenue": 1419,
					"costs": 741,
					"income": 678,
					"dashed": 8620
				},
				{
					"date": "2014-01-13",
					"revenue": 4321,
					"costs": 1428,
					"income": 2893,
					"dashed": 5730
				},
				{
					"date": "2014-01-14",
					"revenue": 1575,
					"costs": 1396,
					"income": 179,
					"dashed": 4778
				},
				{
					"date": "2014-01-15",
					"revenue": 523,
					"costs": 1482,
					"income": -959,
					"dashed": 5581
				},
				{
					"date": "2014-01-16",
					"revenue": 2093,
					"costs": 730,
					"income": 1363,
					"dashed": 4562
				},
				{
					"date": "2014-01-17",
					"revenue": 1681,
					"costs": 353,
					"income": 1328,
					"dashed": 6335
				},
				{
					"date": "2014-01-18",
					"revenue": 3925,
					"costs": 961,
					"income": 2964,
					"dashed": 1976
				},
				{
					"date": "2014-01-19",
					"revenue": 4458,
					"costs": 1004,
					"income": 3454,
					"dashed": 1610
				},
				{
					"date": "2014-01-20",
					"revenue": 3002,
					"costs": 493,
					"income": 2509,
					"dashed": 1565
				},
				{
					"date": "2014-01-21",
					"revenue": 2176,
					"costs": 659,
					"income": 1517,
					"dashed": 6315
				},
				{
					"date": "2014-01-22",
					"revenue": 2763,
					"costs": 444,
					"income": 2319,
					"dashed": 4995
				},
				{
					"date": "2014-01-23",
					"revenue": 3871,
					"costs": 2006,
					"income": 1865,
					"dashed": 9633
				},
				{
					"date": "2014-01-24",
					"revenue": 2740,
					"costs": 1052,
					"income": 1688,
					"dashed": 9702
				},
				{
					"date": "2014-01-25",
					"revenue": 4735,
					"costs": 2312,
					"income": 2423,
					"dashed": 3952
				},
				{
					"date": "2014-01-26",
					"revenue": 3082,
					"costs": 1229,
					"income": 1853,
					"dashed": 9761
				},
				{
					"date": "2014-01-27",
					"revenue": 2374,
					"costs": 244,
					"income": 2130,
					"dashed": 8020
				},
				{
					"date": "2014-01-28",
					"revenue": 2168,
					"costs": 1358,
					"income": 810,
					"dashed": 1447
				},
				{
					"date": "2014-01-29",
					"revenue": 3363,
					"costs": 976,
					"income": 2387,
					"dashed": 2568
				},
				{
					"date": "2014-01-30",
					"revenue": 4465,
					"costs": 751,
					"income": 3714,
					"dashed": 2608
				},
				{
					"date": "2014-01-31",
					"revenue": 4422,
					"costs": 1139,
					"income": 3283,
					"dashed": 4640
				},
				{
					"date": "2014-02-01",
					"revenue": 797,
					"costs": 101,
					"income": 696,
					"dashed": 6379
				},
				{
					"date": "2014-02-02",
					"revenue": 221,
					"costs": 2284,
					"income": -2063,
					"dashed": 3477
				},
				{
					"date": "2014-02-03",
					"revenue": 412,
					"costs": 1312,
					"income": -900,
					"dashed": 7269
				},
				{
					"date": "2014-02-04",
					"revenue": 4002,
					"costs": 1926,
					"income": 2076,
					"dashed": 1031
				},
				{
					"date": "2014-02-05",
					"revenue": 4407,
					"costs": 1803,
					"income": 2604,
					"dashed": 510
				},
				{
					"date": "2014-02-06",
					"revenue": 4664,
					"costs": 2349,
					"income": 2315,
					"dashed": 2473
				},
				{
					"date": "2014-02-07",
					"revenue": 3057,
					"costs": 2467,
					"income": 590,
					"dashed": 1298
				},
				{
					"date": "2014-02-08",
					"revenue": 865,
					"costs": 1836,
					"income": -971,
					"dashed": 5659
				},
				{
					"date": "2014-02-09",
					"revenue": 1663,
					"costs": 993,
					"income": 670,
					"dashed": 943
				},
				{
					"date": "2014-02-10",
					"revenue": 4458,
					"costs": 1432,
					"income": 3026,
					"dashed": 9764
				},
				{
					"date": "2014-02-11",
					"revenue": 4421,
					"costs": 2235,
					"income": 2186,
					"dashed": 6007
				},
				{
					"date": "2014-02-12",
					"revenue": 4039,
					"costs": 2369,
					"income": 1670,
					"dashed": 3952
				},
				{
					"date": "2014-02-13",
					"revenue": 591,
					"costs": 493,
					"income": 98,
					"dashed": 5562
				},
				{
					"date": "2014-02-14",
					"revenue": 3635,
					"costs": 48,
					"income": 3587,
					"dashed": 5958
				},
				{
					"date": "2014-02-15",
					"revenue": 503,
					"costs": 239,
					"income": 264,
					"dashed": 358
				},
				{
					"date": "2014-02-16",
					"revenue": 2485,
					"costs": 972,
					"income": 1513,
					"dashed": 5499
				},
				{
					"date": "2014-02-17",
					"revenue": 1978,
					"costs": 1075,
					"income": 903,
					"dashed": 2436
				},
				{
					"date": "2014-02-18",
					"revenue": 4399,
					"costs": 2262,
					"income": 2137,
					"dashed": 1856
				},
				{
					"date": "2014-02-19",
					"revenue": 1469,
					"costs": 1523,
					"income": -54,
					"dashed": 6902
				},
				{
					"date": "2014-02-20",
					"revenue": 3361,
					"costs": 330,
					"income": 3031,
					"dashed": 2823
				},
				{
					"date": "2014-02-21",
					"revenue": 2487,
					"costs": 1476,
					"income": 1011,
					"dashed": 8635
				},
				{
					"date": "2014-02-22",
					"revenue": 575,
					"costs": 914,
					"income": -339,
					"dashed": 776
				},
				{
					"date": "2014-02-23",
					"revenue": 3083,
					"costs": 2046,
					"income": 1037,
					"dashed": 9008
				},
				{
					"date": "2014-02-24",
					"revenue": 855,
					"costs": 1412,
					"income": -557,
					"dashed": 374
				},
				{
					"date": "2014-02-25",
					"revenue": 3996,
					"costs": 1527,
					"income": 2469,
					"dashed": 7505
				},
				{
					"date": "2014-02-26",
					"revenue": 3428,
					"costs": 766,
					"income": 2662,
					"dashed": 2680
				},
				{
					"date": "2014-02-27",
					"revenue": 4751,
					"costs": 2131,
					"income": 2620,
					"dashed": 3083
				},
				{
					"date": "2014-02-28",
					"revenue": 1706,
					"costs": 1647,
					"income": 59,
					"dashed": 8174
				},
				{
					"date": "2014-03-01",
					"revenue": 4853,
					"costs": 92,
					"income": 4761,
					"dashed": 7450
				},
				{
					"date": "2014-03-02",
					"revenue": 4856,
					"costs": 767,
					"income": 4089,
					"dashed": 5105
				},
				{
					"date": "2014-03-03",
					"revenue": 4441,
					"costs": 1445,
					"income": 2996,
					"dashed": 6846
				},
				{
					"date": "2014-03-04",
					"revenue": 2990,
					"costs": 591,
					"income": 2399,
					"dashed": 8790
				},
				{
					"date": "2014-03-05",
					"revenue": 2309,
					"costs": 2263,
					"income": 46,
					"dashed": 3573
				},
				{
					"date": "2014-03-06",
					"revenue": 3564,
					"costs": 573,
					"income": 2991,
					"dashed": 8374
				},
				{
					"date": "2014-03-07",
					"revenue": 4623,
					"costs": 1463,
					"income": 3160,
					"dashed": 1079
				},
				{
					"date": "2014-03-08",
					"revenue": 4141,
					"costs": 1321,
					"income": 2820,
					"dashed": 4566
				},
				{
					"date": "2014-03-09",
					"revenue": 1651,
					"costs": 1871,
					"income": -220,
					"dashed": 3839
				},
				{
					"date": "2014-03-10",
					"revenue": 2458,
					"costs": 242,
					"income": 2216,
					"dashed": 2835
				},
				{
					"date": "2014-03-11",
					"revenue": 2550,
					"costs": 731,
					"income": 1819,
					"dashed": 4691
				},
				{
					"date": "2014-03-12",
					"revenue": 2795,
					"costs": 584,
					"income": 2211,
					"dashed": 6175
				},
				{
					"date": "2014-03-13",
					"revenue": 3506,
					"costs": 380,
					"income": 3126,
					"dashed": 7159
				},
				{
					"date": "2014-03-14",
					"revenue": 1699,
					"costs": 428,
					"income": 1271,
					"dashed": 9150
				},
				{
					"date": "2014-03-15",
					"revenue": 188,
					"costs": 1430,
					"income": -1242,
					"dashed": 1296
				},
				{
					"date": "2014-03-16",
					"revenue": 158,
					"costs": 835,
					"income": -677,
					"dashed": 1949
				},
				{
					"date": "2014-03-17",
					"revenue": 516,
					"costs": 2173,
					"income": -1657,
					"dashed": 6999
				},
				{
					"date": "2014-03-18",
					"revenue": 1769,
					"costs": 2074,
					"income": -305,
					"dashed": 2283
				},
				{
					"date": "2014-03-19",
					"revenue": 3042,
					"costs": 1934,
					"income": 1108,
					"dashed": 8095
				},
				{
					"date": "2014-03-20",
					"revenue": 1204,
					"costs": 1166,
					"income": 38,
					"dashed": 8771
				},
				{
					"date": "2014-03-21",
					"revenue": 4330,
					"costs": 20,
					"income": 4310,
					"dashed": 889
				},
				{
					"date": "2014-03-22",
					"revenue": 772,
					"costs": 112,
					"income": 660,
					"dashed": 3131
				},
				{
					"date": "2014-03-23",
					"revenue": 3853,
					"costs": 2140,
					"income": 1713,
					"dashed": 5781
				},
				{
					"date": "2014-03-24",
					"revenue": 3324,
					"costs": 1843,
					"income": 1481,
					"dashed": 8743
				},
				{
					"date": "2014-03-25",
					"revenue": 3165,
					"costs": 1487,
					"income": 1678,
					"dashed": 5554
				},
				{
					"date": "2014-03-26",
					"revenue": 4756,
					"costs": 2091,
					"income": 2665,
					"dashed": 644
				},
				{
					"date": "2014-03-27",
					"revenue": 2311,
					"costs": 437,
					"income": 1874,
					"dashed": 286
				},
				{
					"date": "2014-03-28",
					"revenue": 2511,
					"costs": 1323,
					"income": 1188,
					"dashed": 9027
				},
				{
					"date": "2014-03-29",
					"revenue": 1426,
					"costs": 1144,
					"income": 282,
					"dashed": 7
				},
				{
					"date": "2014-03-30",
					"revenue": 3596,
					"costs": 564,
					"income": 3032,
					"dashed": 4450
				},
				{
					"date": "2014-03-31",
					"revenue": 3953,
					"costs": 1711,
					"income": 2242,
					"dashed": 3910
				},
				{
					"date": "2014-04-01",
					"revenue": 3176,
					"costs": 1525,
					"income": 1651,
					"dashed": 7195
				},
				{
					"date": "2014-04-02",
					"revenue": 4786,
					"costs": 824,
					"income": 3962,
					"dashed": 4465
				},
				{
					"date": "2014-04-03",
					"revenue": 417,
					"costs": 1595,
					"income": -1178,
					"dashed": 3730
				},
				{
					"date": "2014-04-04",
					"revenue": 3202,
					"costs": 2232,
					"income": 970,
					"dashed": 9995
				},
				{
					"date": "2014-04-05",
					"revenue": 833,
					"costs": 2173,
					"income": -1340,
					"dashed": 1780
				},
				{
					"date": "2014-04-06",
					"revenue": 1055,
					"costs": 2310,
					"income": -1255,
					"dashed": 3323
				},
				{
					"date": "2014-04-07",
					"revenue": 2251,
					"costs": 1743,
					"income": 508,
					"dashed": 5952
				},
				{
					"date": "2014-04-08",
					"revenue": 190,
					"costs": 2454,
					"income": -2264,
					"dashed": 2579
				},
				{
					"date": "2014-04-09",
					"revenue": 2615,
					"costs": 1591,
					"income": 1024,
					"dashed": 8140
				},
				{
					"date": "2014-04-10",
					"revenue": 4070,
					"costs": 538,
					"income": 3532,
					"dashed": 9050
				},
				{
					"date": "2014-04-11",
					"revenue": 985,
					"costs": 782,
					"income": 203,
					"dashed": 5413
				},
				{
					"date": "2014-04-12",
					"revenue": 4004,
					"costs": 1071,
					"income": 2933,
					"dashed": 9431
				},
				{
					"date": "2014-04-13",
					"revenue": 4360,
					"costs": 25,
					"income": 4335,
					"dashed": 9327
				},
				{
					"date": "2014-04-14",
					"revenue": 1640,
					"costs": 628,
					"income": 1012,
					"dashed": 1480
				},
				{
					"date": "2014-04-15",
					"revenue": 1742,
					"costs": 1206,
					"income": 536,
					"dashed": 9115
				},
				{
					"date": "2014-04-16",
					"revenue": 4905,
					"costs": 2099,
					"income": 2806,
					"dashed": 4785
				},
				{
					"date": "2014-04-17",
					"revenue": 2569,
					"costs": 1470,
					"income": 1099,
					"dashed": 7490
				},
				{
					"date": "2014-04-18",
					"revenue": 842,
					"costs": 2332,
					"income": -1490,
					"dashed": 3733
				},
				{
					"date": "2014-04-19",
					"revenue": 2332,
					"costs": 533,
					"income": 1799,
					"dashed": 429
				},
				{
					"date": "2014-04-20",
					"revenue": 380,
					"costs": 241,
					"income": 139,
					"dashed": 6214
				},
				{
					"date": "2014-04-21",
					"revenue": 3598,
					"costs": 2213,
					"income": 1385,
					"dashed": 6863
				},
				{
					"date": "2014-04-22",
					"revenue": 2240,
					"costs": 2297,
					"income": -57,
					"dashed": 915
				},
				{
					"date": "2014-04-23",
					"revenue": 4232,
					"costs": 2406,
					"income": 1826,
					"dashed": 1219
				},
				{
					"date": "2014-04-24",
					"revenue": 2615,
					"costs": 2180,
					"income": 435,
					"dashed": 9037
				},
				{
					"date": "2014-04-25",
					"revenue": 3448,
					"costs": 2465,
					"income": 983,
					"dashed": 6075
				},
				{
					"date": "2014-04-26",
					"revenue": 4530,
					"costs": 2193,
					"income": 2337,
					"dashed": 6423
				},
				{
					"date": "2014-04-27",
					"revenue": 4194,
					"costs": 716,
					"income": 3478,
					"dashed": 9983
				},
				{
					"date": "2014-04-28",
					"revenue": 96,
					"costs": 423,
					"income": -327,
					"dashed": 6826
				},
				{
					"date": "2014-04-29",
					"revenue": 4358,
					"costs": 912,
					"income": 3446,
					"dashed": 1479
				},
				{
					"date": "2014-04-30",
					"revenue": 3930,
					"costs": 1992,
					"income": 1938,
					"dashed": 9882
				},
				{
					"date": "2014-05-01",
					"revenue": 4802,
					"costs": 2128,
					"income": 2674,
					"dashed": 5589
				},
				{
					"date": "2014-05-02",
					"revenue": 1850,
					"costs": 862,
					"income": 988,
					"dashed": 1525
				},
				{
					"date": "2014-05-03",
					"revenue": 3110,
					"costs": 1107,
					"income": 2003,
					"dashed": 9486
				},
				{
					"date": "2014-05-04",
					"revenue": 2198,
					"costs": 2082,
					"income": 116,
					"dashed": 851
				},
				{
					"date": "2014-05-05",
					"revenue": 1930,
					"costs": 1865,
					"income": 65,
					"dashed": 4906
				},
				{
					"date": "2014-05-06",
					"revenue": 291,
					"costs": 1253,
					"income": -962,
					"dashed": 7060
				},
				{
					"date": "2014-05-07",
					"revenue": 1935,
					"costs": 31,
					"income": 1904,
					"dashed": 9193
				},
				{
					"date": "2014-05-08",
					"revenue": 2848,
					"costs": 2414,
					"income": 434,
					"dashed": 8282
				},
				{
					"date": "2014-05-09",
					"revenue": 4672,
					"costs": 1807,
					"income": 2865,
					"dashed": 9648
				},
				{
					"date": "2014-05-10",
					"revenue": 1743,
					"costs": 797,
					"income": 946,
					"dashed": 65
				},
				{
					"date": "2014-05-11",
					"revenue": 2532,
					"costs": 940,
					"income": 1592,
					"dashed": 7348
				},
				{
					"date": "2014-05-12",
					"revenue": 4098,
					"costs": 2284,
					"income": 1814,
					"dashed": 6733
				},
				{
					"date": "2014-05-13",
					"revenue": 4567,
					"costs": 2236,
					"income": 2331,
					"dashed": 691
				},
				{
					"date": "2014-05-14",
					"revenue": 109,
					"costs": 2328,
					"income": -2219,
					"dashed": 726
				},
				{
					"date": "2014-05-15",
					"revenue": 4303,
					"costs": 1449,
					"income": 2854,
					"dashed": 5052
				},
				{
					"date": "2014-05-16",
					"revenue": 3745,
					"costs": 1140,
					"income": 2605,
					"dashed": 3005
				},
				{
					"date": "2014-05-17",
					"revenue": 2940,
					"costs": 472,
					"income": 2468,
					"dashed": 9008
				},
				{
					"date": "2014-05-18",
					"revenue": 4208,
					"costs": 1260,
					"income": 2948,
					"dashed": 2308
				},
				{
					"date": "2014-05-19",
					"revenue": 3776,
					"costs": 725,
					"income": 3051,
					"dashed": 8398
				},
				{
					"date": "2014-05-20",
					"revenue": 501,
					"costs": 539,
					"income": -38,
					"dashed": 7182
				},
				{
					"date": "2014-05-21",
					"revenue": 1944,
					"costs": 503,
					"income": 1441,
					"dashed": 7459
				},
				{
					"date": "2014-05-22",
					"revenue": 3216,
					"costs": 482,
					"income": 2734,
					"dashed": 6491
				},
				{
					"date": "2014-05-23",
					"revenue": 3428,
					"costs": 493,
					"income": 2935,
					"dashed": 2900
				},
				{
					"date": "2014-05-24",
					"revenue": 1497,
					"costs": 1716,
					"income": -219,
					"dashed": 7982
				},
				{
					"date": "2014-05-25",
					"revenue": 3329,
					"costs": 266,
					"income": 3063,
					"dashed": 7252
				},
				{
					"date": "2014-05-26",
					"revenue": 2890,
					"costs": 1710,
					"income": 1180,
					"dashed": 8924
				},
				{
					"date": "2014-05-27",
					"revenue": 3381,
					"costs": 2208,
					"income": 1173,
					"dashed": 2368
				},
				{
					"date": "2014-05-28",
					"revenue": 2401,
					"costs": 1556,
					"income": 845,
					"dashed": 4067
				},
				{
					"date": "2014-05-29",
					"revenue": 1423,
					"costs": 726,
					"income": 697,
					"dashed": 9387
				},
				{
					"date": "2014-05-30",
					"revenue": 2523,
					"costs": 171,
					"income": 2352,
					"dashed": 332
				},
				{
					"date": "2014-05-31",
					"revenue": 3483,
					"costs": 863,
					"income": 2620,
					"dashed": 7330
				},
				{
					"date": "2014-06-01",
					"revenue": 2575,
					"costs": 2113,
					"income": 462,
					"dashed": 9069
				},
				{
					"date": "2014-06-02",
					"revenue": 604,
					"costs": 1802,
					"income": -1198,
					"dashed": 6452
				},
				{
					"date": "2014-06-03",
					"revenue": 2149,
					"costs": 1851,
					"income": 298,
					"dashed": 3410
				},
				{
					"date": "2014-06-04",
					"revenue": 4335,
					"costs": 1376,
					"income": 2959,
					"dashed": 9354
				},
				{
					"date": "2014-06-05",
					"revenue": 3123,
					"costs": 1069,
					"income": 2054,
					"dashed": 3957
				},
				{
					"date": "2014-06-06",
					"revenue": 2475,
					"costs": 2462,
					"income": 13,
					"dashed": 4780
				},
				{
					"date": "2014-06-07",
					"revenue": 2418,
					"costs": 317,
					"income": 2101,
					"dashed": 6667
				},
				{
					"date": "2014-06-08",
					"revenue": 609,
					"costs": 2015,
					"income": -1406,
					"dashed": 1315
				},
				{
					"date": "2014-06-09",
					"revenue": 1246,
					"costs": 1042,
					"income": 204,
					"dashed": 2856
				},
				{
					"date": "2014-06-10",
					"revenue": 571,
					"costs": 641,
					"income": -70,
					"dashed": 7961
				},
				{
					"date": "2014-06-11",
					"revenue": 878,
					"costs": 206,
					"income": 672,
					"dashed": 5588
				},
				{
					"date": "2014-06-12",
					"revenue": 1740,
					"costs": 1129,
					"income": 611,
					"dashed": 8059
				},
				{
					"date": "2014-06-13",
					"revenue": 4143,
					"costs": 2404,
					"income": 1739,
					"dashed": 8723
				},
				{
					"date": "2014-06-14",
					"revenue": 1549,
					"costs": 1758,
					"income": -209,
					"dashed": 7129
				},
				{
					"date": "2014-06-15",
					"revenue": 1690,
					"costs": 2078,
					"income": -388,
					"dashed": 9956
				},
				{
					"date": "2014-06-16",
					"revenue": 3995,
					"costs": 2145,
					"income": 1850,
					"dashed": 1051
				},
				{
					"date": "2014-06-17",
					"revenue": 4382,
					"costs": 672,
					"income": 3710,
					"dashed": 7998
				},
				{
					"date": "2014-06-18",
					"revenue": 32,
					"costs": 560,
					"income": -528,
					"dashed": 7759
				},
				{
					"date": "2014-06-19",
					"revenue": 4593,
					"costs": 1397,
					"income": 3196,
					"dashed": 4978
				},
				{
					"date": "2014-06-20",
					"revenue": 13,
					"costs": 2479,
					"income": -2466,
					"dashed": 4247
				},
				{
					"date": "2014-06-21",
					"revenue": 3080,
					"costs": 584,
					"income": 2496,
					"dashed": 5
				},
				{
					"date": "2014-06-22",
					"revenue": 2174,
					"costs": 2180,
					"income": -6,
					"dashed": 4214
				},
				{
					"date": "2014-06-23",
					"revenue": 572,
					"costs": 362,
					"income": 210,
					"dashed": 460
				},
				{
					"date": "2014-06-24",
					"revenue": 3191,
					"costs": 2081,
					"income": 1110,
					"dashed": 3813
				},
				{
					"date": "2014-06-25",
					"revenue": 1570,
					"costs": 1419,
					"income": 151,
					"dashed": 2496
				},
				{
					"date": "2014-06-26",
					"revenue": 2104,
					"costs": 323,
					"income": 1781,
					"dashed": 5596
				},
				{
					"date": "2014-06-27",
					"revenue": 4694,
					"costs": 717,
					"income": 3977,
					"dashed": 4784
				},
				{
					"date": "2014-06-28",
					"revenue": 1719,
					"costs": 2010,
					"income": -291,
					"dashed": 2923
				},
				{
					"date": "2014-06-29",
					"revenue": 2041,
					"costs": 1391,
					"income": 650,
					"dashed": 7431
				},
				{
					"date": "2014-06-30",
					"revenue": 527,
					"costs": 1307,
					"income": -780,
					"dashed": 4074
				},
				{
					"date": "2014-07-01",
					"revenue": 1916,
					"costs": 2486,
					"income": -570,
					"dashed": 9763
				},
				{
					"date": "2014-07-02",
					"revenue": 4336,
					"costs": 1220,
					"income": 3116,
					"dashed": 5549
				},
				{
					"date": "2014-07-03",
					"revenue": 4208,
					"costs": 1276,
					"income": 2932,
					"dashed": 6095
				},
				{
					"date": "2014-07-04",
					"revenue": 4206,
					"costs": 2458,
					"income": 1748,
					"dashed": 4026
				},
				{
					"date": "2014-07-05",
					"revenue": 3974,
					"costs": 1738,
					"income": 2236,
					"dashed": 7334
				},
				{
					"date": "2014-07-06",
					"revenue": 1112,
					"costs": 295,
					"income": 817,
					"dashed": 7225
				},
				{
					"date": "2014-07-07",
					"revenue": 2279,
					"costs": 345,
					"income": 1934,
					"dashed": 9051
				},
				{
					"date": "2014-07-08",
					"revenue": 916,
					"costs": 2344,
					"income": -1428,
					"dashed": 3660
				},
				{
					"date": "2014-07-09",
					"revenue": 3770,
					"costs": 1591,
					"income": 2179,
					"dashed": 1924
				},
				{
					"date": "2014-07-10",
					"revenue": 1635,
					"costs": 783,
					"income": 852,
					"dashed": 2477
				},
				{
					"date": "2014-07-11",
					"revenue": 4923,
					"costs": 672,
					"income": 4251,
					"dashed": 9553
				},
				{
					"date": "2014-07-12",
					"revenue": 1963,
					"costs": 1833,
					"income": 130,
					"dashed": 9919
				},
				{
					"date": "2014-07-13",
					"revenue": 3543,
					"costs": 2260,
					"income": 1283,
					"dashed": 8570
				},
				{
					"date": "2014-07-14",
					"revenue": 3677,
					"costs": 313,
					"income": 3364,
					"dashed": 3712
				},
				{
					"date": "2014-07-15",
					"revenue": 1650,
					"costs": 176,
					"income": 1474,
					"dashed": 322
				},
				{
					"date": "2014-07-16",
					"revenue": 4760,
					"costs": 314,
					"income": 4446,
					"dashed": 8118
				},
				{
					"date": "2014-07-17",
					"revenue": 2194,
					"costs": 2132,
					"income": 62,
					"dashed": 554
				},
				{
					"date": "2014-07-18",
					"revenue": 4629,
					"costs": 544,
					"income": 4085,
					"dashed": 9034
				},
				{
					"date": "2014-07-19",
					"revenue": 4702,
					"costs": 439,
					"income": 4263,
					"dashed": 7172
				},
				{
					"date": "2014-07-20",
					"revenue": 4816,
					"costs": 138,
					"income": 4678,
					"dashed": 3242
				},
				{
					"date": "2014-07-21",
					"revenue": 1259,
					"costs": 2478,
					"income": -1219,
					"dashed": 7665
				},
				{
					"date": "2014-07-22",
					"revenue": 1947,
					"costs": 91,
					"income": 1856,
					"dashed": 8381
				},
				{
					"date": "2014-07-23",
					"revenue": 2921,
					"costs": 828,
					"income": 2093,
					"dashed": 7873
				},
				{
					"date": "2014-07-24",
					"revenue": 4958,
					"costs": 2462,
					"income": 2496,
					"dashed": 6109
				},
				{
					"date": "2014-07-25",
					"revenue": 1720,
					"costs": 1623,
					"income": 97,
					"dashed": 8819
				},
				{
					"date": "2014-07-26",
					"revenue": 4272,
					"costs": 437,
					"income": 3835,
					"dashed": 4096
				},
				{
					"date": "2014-07-27",
					"revenue": 3641,
					"costs": 2344,
					"income": 1297,
					"dashed": 9577
				},
				{
					"date": "2014-07-28",
					"revenue": 4571,
					"costs": 654,
					"income": 3917,
					"dashed": 8119
				},
				{
					"date": "2014-07-29",
					"revenue": 4714,
					"costs": 2251,
					"income": 2463,
					"dashed": 4181
				},
				{
					"date": "2014-07-30",
					"revenue": 1731,
					"costs": 802,
					"income": 929,
					"dashed": 3577
				},
				{
					"date": "2014-07-31",
					"revenue": 3682,
					"costs": 128,
					"income": 3554,
					"dashed": 2910
				},
				{
					"date": "2014-08-01",
					"revenue": 259,
					"costs": 1574,
					"income": -1315,
					"dashed": 4107
				},
				{
					"date": "2014-08-02",
					"revenue": 4094,
					"costs": 2444,
					"income": 1650,
					"dashed": 2883
				},
				{
					"date": "2014-08-03",
					"revenue": 1420,
					"costs": 1724,
					"income": -304,
					"dashed": 7104
				},
				{
					"date": "2014-08-04",
					"revenue": 4939,
					"costs": 2038,
					"income": 2901,
					"dashed": 9687
				},
				{
					"date": "2014-08-05",
					"revenue": 3809,
					"costs": 1969,
					"income": 1840,
					"dashed": 5870
				},
				{
					"date": "2014-08-06",
					"revenue": 3512,
					"costs": 1451,
					"income": 2061,
					"dashed": 4681
				},
				{
					"date": "2014-08-07",
					"revenue": 3695,
					"costs": 837,
					"income": 2858,
					"dashed": 4969
				},
				{
					"date": "2014-08-08",
					"revenue": 2638,
					"costs": 1023,
					"income": 1615,
					"dashed": 133
				},
				{
					"date": "2014-08-09",
					"revenue": 4425,
					"costs": 448,
					"income": 3977,
					"dashed": 2727
				},
				{
					"date": "2014-08-10",
					"revenue": 2805,
					"costs": 1862,
					"income": 943,
					"dashed": 3443
				},
				{
					"date": "2014-08-11",
					"revenue": 427,
					"costs": 2172,
					"income": -1745,
					"dashed": 8399
				},
				{
					"date": "2014-08-12",
					"revenue": 3791,
					"costs": 1787,
					"income": 2004,
					"dashed": 3421
				},
				{
					"date": "2014-08-13",
					"revenue": 1453,
					"costs": 1617,
					"income": -164,
					"dashed": 1017
				},
				{
					"date": "2014-08-14",
					"revenue": 671,
					"costs": 2321,
					"income": -1650,
					"dashed": 3321
				},
				{
					"date": "2014-08-15",
					"revenue": 2853,
					"costs": 1221,
					"income": 1632,
					"dashed": 263
				},
				{
					"date": "2014-08-16",
					"revenue": 4705,
					"costs": 1339,
					"income": 3366,
					"dashed": 8853
				},
				{
					"date": "2014-08-17",
					"revenue": 2182,
					"costs": 64,
					"income": 2118,
					"dashed": 1777
				},
				{
					"date": "2014-08-18",
					"revenue": 715,
					"costs": 1510,
					"income": -795,
					"dashed": 2666
				},
				{
					"date": "2014-08-19",
					"revenue": 389,
					"costs": 2235,
					"income": -1846,
					"dashed": 8791
				},
				{
					"date": "2014-08-20",
					"revenue": 490,
					"costs": 628,
					"income": -138,
					"dashed": 3944
				},
				{
					"date": "2014-08-21",
					"revenue": 4438,
					"costs": 795,
					"income": 3643,
					"dashed": 9156
				},
				{
					"date": "2014-08-22",
					"revenue": 320,
					"costs": 535,
					"income": -215,
					"dashed": 5232
				},
				{
					"date": "2014-08-23",
					"revenue": 177,
					"costs": 1898,
					"income": -1721,
					"dashed": 3998
				},
				{
					"date": "2014-08-24",
					"revenue": 3184,
					"costs": 1016,
					"income": 2168,
					"dashed": 322
				},
				{
					"date": "2014-08-25",
					"revenue": 577,
					"costs": 350,
					"income": 227,
					"dashed": 6485
				},
				{
					"date": "2014-08-26",
					"revenue": 277,
					"costs": 1013,
					"income": -736,
					"dashed": 3394
				},
				{
					"date": "2014-08-27",
					"revenue": 414,
					"costs": 1136,
					"income": -722,
					"dashed": 9021
				},
				{
					"date": "2014-08-28",
					"revenue": 1505,
					"costs": 1285,
					"income": 220,
					"dashed": 4713
				},
				{
					"date": "2014-08-29",
					"revenue": 921,
					"costs": 816,
					"income": 105,
					"dashed": 6951
				},
				{
					"date": "2014-08-30",
					"revenue": 646,
					"costs": 574,
					"income": 72,
					"dashed": 2287
				},
				{
					"date": "2014-08-31",
					"revenue": 3473,
					"costs": 1899,
					"income": 1574,
					"dashed": 3700
				},
				{
					"date": "2014-09-01",
					"revenue": 3082,
					"costs": 2074,
					"income": 1008,
					"dashed": 631
				},
				{
					"date": "2014-09-02",
					"revenue": 3200,
					"costs": 120,
					"income": 3080,
					"dashed": 7341
				},
				{
					"date": "2014-09-03",
					"revenue": 3822,
					"costs": 755,
					"income": 3067,
					"dashed": 5715
				},
				{
					"date": "2014-09-04",
					"revenue": 4717,
					"costs": 886,
					"income": 3831,
					"dashed": 6248
				},
				{
					"date": "2014-09-05",
					"revenue": 3925,
					"costs": 2178,
					"income": 1747,
					"dashed": 5740
				},
				{
					"date": "2014-09-06",
					"revenue": 3875,
					"costs": 642,
					"income": 3233,
					"dashed": 2951
				},
				{
					"date": "2014-09-07",
					"revenue": 3743,
					"costs": 2076,
					"income": 1667,
					"dashed": 507
				},
				{
					"date": "2014-09-08",
					"revenue": 500,
					"costs": 888,
					"income": -388,
					"dashed": 991
				},
				{
					"date": "2014-09-09",
					"revenue": 4644,
					"costs": 1727,
					"income": 2917,
					"dashed": 5579
				},
				{
					"date": "2014-09-10",
					"revenue": 3096,
					"costs": 382,
					"income": 2714,
					"dashed": 217
				},
				{
					"date": "2014-09-11",
					"revenue": 3683,
					"costs": 232,
					"income": 3451,
					"dashed": 4700
				},
				{
					"date": "2014-09-12",
					"revenue": 1640,
					"costs": 291,
					"income": 1349,
					"dashed": 1498
				},
				{
					"date": "2014-09-13",
					"revenue": 4763,
					"costs": 1117,
					"income": 3646,
					"dashed": 3025
				},
				{
					"date": "2014-09-14",
					"revenue": 1575,
					"costs": 371,
					"income": 1204,
					"dashed": 1385
				},
				{
					"date": "2014-09-15",
					"revenue": 2631,
					"costs": 1568,
					"income": 1063,
					"dashed": 9678
				},
				{
					"date": "2014-09-16",
					"revenue": 4281,
					"costs": 2027,
					"income": 2254,
					"dashed": 7144
				},
				{
					"date": "2014-09-17",
					"revenue": 4647,
					"costs": 1173,
					"income": 3474,
					"dashed": 2098
				},
				{
					"date": "2014-09-18",
					"revenue": 1822,
					"costs": 1791,
					"income": 31,
					"dashed": 4297
				},
				{
					"date": "2014-09-19",
					"revenue": 761,
					"costs": 1487,
					"income": -726,
					"dashed": 9044
				},
				{
					"date": "2014-09-20",
					"revenue": 2685,
					"costs": 729,
					"income": 1956,
					"dashed": 616
				},
				{
					"date": "2014-09-21",
					"revenue": 1661,
					"costs": 773,
					"income": 888,
					"dashed": 3349
				},
				{
					"date": "2014-09-22",
					"revenue": 2743,
					"costs": 17,
					"income": 2726,
					"dashed": 8259
				},
				{
					"date": "2014-09-23",
					"revenue": 4608,
					"costs": 1345,
					"income": 3263,
					"dashed": 1042
				},
				{
					"date": "2014-09-24",
					"revenue": 1267,
					"costs": 1301,
					"income": -34,
					"dashed": 3088
				},
				{
					"date": "2014-09-25",
					"revenue": 2127,
					"costs": 1726,
					"income": 401,
					"dashed": 2334
				},
				{
					"date": "2014-09-26",
					"revenue": 2271,
					"costs": 2281,
					"income": -10,
					"dashed": 5247
				},
				{
					"date": "2014-09-27",
					"revenue": 1383,
					"costs": 1203,
					"income": 180,
					"dashed": 1510
				},
				{
					"date": "2014-09-28",
					"revenue": 2965,
					"costs": 1175,
					"income": 1790,
					"dashed": 4496
				},
				{
					"date": "2014-09-29",
					"revenue": 3334,
					"costs": 271,
					"income": 3063,
					"dashed": 2382
				},
				{
					"date": "2014-09-30",
					"revenue": 278,
					"costs": 421,
					"income": -143,
					"dashed": 3750
				},
				{
					"date": "2014-10-01",
					"revenue": 1989,
					"costs": 1165,
					"income": 824,
					"dashed": 7319
				},
				{
					"date": "2014-10-02",
					"revenue": 3273,
					"costs": 1817,
					"income": 1456,
					"dashed": 3268
				},
				{
					"date": "2014-10-03",
					"revenue": 4728,
					"costs": 1219,
					"income": 3509,
					"dashed": 6640
				},
				{
					"date": "2014-10-04",
					"revenue": 2760,
					"costs": 1241,
					"income": 1519,
					"dashed": 2925
				},
				{
					"date": "2014-10-05",
					"revenue": 4124,
					"costs": 977,
					"income": 3147,
					"dashed": 4482
				},
				{
					"date": "2014-10-06",
					"revenue": 1752,
					"costs": 1044,
					"income": 708,
					"dashed": 4413
				},
				{
					"date": "2014-10-07",
					"revenue": 205,
					"costs": 941,
					"income": -736,
					"dashed": 3880
				},
				{
					"date": "2014-10-08",
					"revenue": 4007,
					"costs": 337,
					"income": 3670,
					"dashed": 7805
				},
				{
					"date": "2014-10-09",
					"revenue": 3960,
					"costs": 1273,
					"income": 2687,
					"dashed": 8150
				},
				{
					"date": "2014-10-10",
					"revenue": 3874,
					"costs": 1121,
					"income": 2753,
					"dashed": 543
				},
				{
					"date": "2014-10-11",
					"revenue": 4501,
					"costs": 361,
					"income": 4140,
					"dashed": 1586
				},
				{
					"date": "2014-10-12",
					"revenue": 159,
					"costs": 54,
					"income": 105,
					"dashed": 8787
				},
				{
					"date": "2014-10-13",
					"revenue": 3774,
					"costs": 1413,
					"income": 2361,
					"dashed": 1728
				},
				{
					"date": "2014-10-14",
					"revenue": 1197,
					"costs": 1108,
					"income": 89,
					"dashed": 279
				},
				{
					"date": "2014-10-15",
					"revenue": 530,
					"costs": 1703,
					"income": -1173,
					"dashed": 9112
				},
				{
					"date": "2014-10-16",
					"revenue": 1886,
					"costs": 1791,
					"income": 95,
					"dashed": 7297
				},
				{
					"date": "2014-10-17",
					"revenue": 1837,
					"costs": 689,
					"income": 1148,
					"dashed": 45
				},
				{
					"date": "2014-10-18",
					"revenue": 3739,
					"costs": 1873,
					"income": 1866,
					"dashed": 1318
				},
				{
					"date": "2014-10-19",
					"revenue": 4353,
					"costs": 959,
					"income": 3394,
					"dashed": 2017
				},
				{
					"date": "2014-10-20",
					"revenue": 1516,
					"costs": 1917,
					"income": -401,
					"dashed": 780
				},
				{
					"date": "2014-10-21",
					"revenue": 3461,
					"costs": 1665,
					"income": 1796,
					"dashed": 4965
				},
				{
					"date": "2014-10-22",
					"revenue": 4079,
					"costs": 1411,
					"income": 2668,
					"dashed": 3458
				},
				{
					"date": "2014-10-23",
					"revenue": 3629,
					"costs": 1727,
					"income": 1902,
					"dashed": 81
				},
				{
					"date": "2014-10-24",
					"revenue": 3094,
					"costs": 520,
					"income": 2574,
					"dashed": 9382
				},
				{
					"date": "2014-10-25",
					"revenue": 3080,
					"costs": 1953,
					"income": 1127,
					"dashed": 618
				},
				{
					"date": "2014-10-26",
					"revenue": 2693,
					"costs": 1635,
					"income": 1058,
					"dashed": 9149
				},
				{
					"date": "2014-10-27",
					"revenue": 4376,
					"costs": 128,
					"income": 4248,
					"dashed": 3518
				},
				{
					"date": "2014-10-28",
					"revenue": 4135,
					"costs": 1465,
					"income": 2670,
					"dashed": 1582
				},
				{
					"date": "2014-10-29",
					"revenue": 3110,
					"costs": 333,
					"income": 2777,
					"dashed": 3753
				},
				{
					"date": "2014-10-30",
					"revenue": 1185,
					"costs": 2111,
					"income": -926,
					"dashed": 6724
				},
				{
					"date": "2014-10-31",
					"revenue": 693,
					"costs": 1947,
					"income": -1254,
					"dashed": 1931
				},
				{
					"date": "2014-11-01",
					"revenue": 1420,
					"costs": 1240,
					"income": 180,
					"dashed": 5125
				},
				{
					"date": "2014-11-02",
					"revenue": 579,
					"costs": 249,
					"income": 330,
					"dashed": 2232
				},
				{
					"date": "2014-11-03",
					"revenue": 1062,
					"costs": 2396,
					"income": -1334,
					"dashed": 3393
				},
				{
					"date": "2014-11-04",
					"revenue": 3259,
					"costs": 906,
					"income": 2353,
					"dashed": 2581
				},
				{
					"date": "2014-11-05",
					"revenue": 1241,
					"costs": 1560,
					"income": -319,
					"dashed": 164
				},
				{
					"date": "2014-11-06",
					"revenue": 975,
					"costs": 950,
					"income": 25,
					"dashed": 1583
				},
				{
					"date": "2014-11-07",
					"revenue": 2537,
					"costs": 78,
					"income": 2459,
					"dashed": 3674
				},
				{
					"date": "2014-11-08",
					"revenue": 350,
					"costs": 2447,
					"income": -2097,
					"dashed": 5946
				},
				{
					"date": "2014-11-09",
					"revenue": 2692,
					"costs": 215,
					"income": 2477,
					"dashed": 953
				},
				{
					"date": "2014-11-10",
					"revenue": 4231,
					"costs": 1665,
					"income": 2566,
					"dashed": 7036
				},
				{
					"date": "2014-11-11",
					"revenue": 4149,
					"costs": 2074,
					"income": 2075,
					"dashed": 7015
				},
				{
					"date": "2014-11-12",
					"revenue": 3228,
					"costs": 2229,
					"income": 999,
					"dashed": 4754
				},
				{
					"date": "2014-11-13",
					"revenue": 3067,
					"costs": 92,
					"income": 2975,
					"dashed": 6024
				},
				{
					"date": "2014-11-14",
					"revenue": 96,
					"costs": 559,
					"income": -463,
					"dashed": 7169
				},
				{
					"date": "2014-11-15",
					"revenue": 1535,
					"costs": 1762,
					"income": -227,
					"dashed": 2621
				},
				{
					"date": "2014-11-16",
					"revenue": 1968,
					"costs": 2159,
					"income": -191,
					"dashed": 4929
				},
				{
					"date": "2014-11-17",
					"revenue": 5,
					"costs": 167,
					"income": -162,
					"dashed": 7554
				},
				{
					"date": "2014-11-18",
					"revenue": 1770,
					"costs": 1101,
					"income": 669,
					"dashed": 2206
				},
				{
					"date": "2014-11-19",
					"revenue": 1121,
					"costs": 913,
					"income": 208,
					"dashed": 4821
				},
				{
					"date": "2014-11-20",
					"revenue": 1741,
					"costs": 732,
					"income": 1009,
					"dashed": 3593
				},
				{
					"date": "2014-11-21",
					"revenue": 3182,
					"costs": 921,
					"income": 2261,
					"dashed": 8375
				},
				{
					"date": "2014-11-22",
					"revenue": 3463,
					"costs": 342,
					"income": 3121,
					"dashed": 3957
				},
				{
					"date": "2014-11-23",
					"revenue": 3201,
					"costs": 881,
					"income": 2320,
					"dashed": 3537
				},
				{
					"date": "2014-11-24",
					"revenue": 3333,
					"costs": 932,
					"income": 2401,
					"dashed": 9479
				},
				{
					"date": "2014-11-25",
					"revenue": 835,
					"costs": 513,
					"income": 322,
					"dashed": 9799
				},
				{
					"date": "2014-11-26",
					"revenue": 1710,
					"costs": 1811,
					"income": -101,
					"dashed": 1560
				},
				{
					"date": "2014-11-27",
					"revenue": 287,
					"costs": 437,
					"income": -150,
					"dashed": 96
				},
				{
					"date": "2014-11-28",
					"revenue": 2499,
					"costs": 2402,
					"income": 97,
					"dashed": 5028
				},
				{
					"date": "2014-11-29",
					"revenue": 3833,
					"costs": 1766,
					"income": 2067,
					"dashed": 1002
				},
				{
					"date": "2014-11-30",
					"revenue": 1627,
					"costs": 1170,
					"income": 457,
					"dashed": 3079
				},
				{
					"date": "2014-12-01",
					"revenue": 4092,
					"costs": 2155,
					"income": 1937,
					"dashed": 298
				},
				{
					"date": "2014-12-02",
					"revenue": 1372,
					"costs": 463,
					"income": 909,
					"dashed": 5615
				},
				{
					"date": "2014-12-03",
					"revenue": 1479,
					"costs": 539,
					"income": 940,
					"dashed": 9070
				},
				{
					"date": "2014-12-04",
					"revenue": 4886,
					"costs": 1,
					"income": 4885,
					"dashed": 1585
				},
				{
					"date": "2014-12-05",
					"revenue": 4198,
					"costs": 45,
					"income": 4153,
					"dashed": 9867
				},
				{
					"date": "2014-12-06",
					"revenue": 4810,
					"costs": 21,
					"income": 4789,
					"dashed": 2183
				},
				{
					"date": "2014-12-07",
					"revenue": 147,
					"costs": 378,
					"income": -231,
					"dashed": 6553
				},
				{
					"date": "2014-12-08",
					"revenue": 4807,
					"costs": 324,
					"income": 4483,
					"dashed": 8464
				},
				{
					"date": "2014-12-09",
					"revenue": 3359,
					"costs": 1424,
					"income": 1935,
					"dashed": 3404
				},
				{
					"date": "2014-12-10",
					"revenue": 2110,
					"costs": 318,
					"income": 1792,
					"dashed": 7860
				},
				{
					"date": "2014-12-11",
					"revenue": 992,
					"costs": 147,
					"income": 845,
					"dashed": 6405
				},
				{
					"date": "2014-12-12",
					"revenue": 4239,
					"costs": 996,
					"income": 3243,
					"dashed": 2167
				},
				{
					"date": "2014-12-13",
					"revenue": 1025,
					"costs": 1915,
					"income": -890,
					"dashed": 5378
				},
				{
					"date": "2014-12-14",
					"revenue": 3124,
					"costs": 863,
					"income": 2261,
					"dashed": 4270
				},
				{
					"date": "2014-12-15",
					"revenue": 3577,
					"costs": 1448,
					"income": 2129,
					"dashed": 8032
				},
				{
					"date": "2014-12-16",
					"revenue": 3913,
					"costs": 915,
					"income": 2998,
					"dashed": 2181
				},
				{
					"date": "2014-12-17",
					"revenue": 4350,
					"costs": 856,
					"income": 3494,
					"dashed": 7719
				},
				{
					"date": "2014-12-18",
					"revenue": 4528,
					"costs": 2424,
					"income": 2104,
					"dashed": 1834
				},
				{
					"date": "2014-12-19",
					"revenue": 3252,
					"costs": 375,
					"income": 2877,
					"dashed": 5577
				},
				{
					"date": "2014-12-20",
					"revenue": 153,
					"costs": 1327,
					"income": -1174,
					"dashed": 6171
				},
				{
					"date": "2014-12-21",
					"revenue": 1193,
					"costs": 618,
					"income": 575,
					"dashed": 7588
				},
				{
					"date": "2014-12-22",
					"revenue": 726,
					"costs": 961,
					"income": -235,
					"dashed": 2806
				},
				{
					"date": "2014-12-23",
					"revenue": 2245,
					"costs": 2342,
					"income": -97,
					"dashed": 2972
				},
				{
					"date": "2014-12-24",
					"revenue": 4389,
					"costs": 170,
					"income": 4219,
					"dashed": 9920
				},
				{
					"date": "2014-12-25",
					"revenue": 4306,
					"costs": 1806,
					"income": 2500,
					"dashed": 7772
				},
				{
					"date": "2014-12-26",
					"revenue": 4276,
					"costs": 2269,
					"income": 2007,
					"dashed": 3633
				},
				{
					"date": "2014-12-27",
					"revenue": 1136,
					"costs": 1270,
					"income": -134,
					"dashed": 717
				},
				{
					"date": "2014-12-28",
					"revenue": 2750,
					"costs": 1129,
					"income": 1621,
					"dashed": 4484
				},
				{
					"date": "2014-12-29",
					"revenue": 485,
					"costs": 1914,
					"income": -1429,
					"dashed": 1296
				},
				{
					"date": "2014-12-30",
					"revenue": 4992,
					"costs": 1582,
					"income": 3410,
					"dashed": 5231
				}]
      },
      theme: "flat",
      seriesDefaults: {
	        area: {
	            line: {
	                style: "smooth"
	            }
	        }
	    },
      dateField: "date",
      series: [{
        type: "area",
        field: "revenue",
        aggregate: "sum", 
        color: "#ea5b19",
      }],
      navigator: {
        categoryAxis: {
          roundToBaseUnit: true
        },
        series: [{
          type: "area",
          field: "revenue",
          aggregate: "sum",
          color: "#ea5b19",
        }]
      },
      valueAxis: {
        title: {
            text: "Production",
            visible: true,
            font: '12px Source Sans Pro, Lato , Open Sans , Helvetica Neue, Arial, sans-serif'
        },
        labels: {
            // format: "{0:p2}"
            format: "{0}"
        },
        majorGridLines: {
            visible: true,
            color: "#eee",
            width: 0.8,
        },
        line: {
            visible: false
        },
        axisCrossingValue: 0
      },
	  categoryAxis: {
            majorGridLines: {
                visible: false
            },
            majorTickType: "none"
      },
    });
} 
$(document).ready(function () {
    fa.LoadData();
    $('#btnRefresh').on('click', function () {
        fa.LoadData();
    });

    setTimeout(function () {
        fa.LoadData();
        pg.chartWindSpeed();
        pg.chartProduction();
    }, 1000);
    app.loading(false);
});