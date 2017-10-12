### Entity包介绍(持续开发中)
```
名词介绍：
Entity：实体，是player,monster,npc甚至是防御塔等游戏实体的基类，内置了aoi，属性，道具
等属性。
Space：空间，本身也是一种Entity，这里定义为Entity的容器，既然叫做空间，自然拥有AOI计算的能力。
AOI：Area Of Interest，即实时关注区域，实质是一种数据结构，用于记录变动的位置和关注信息。
Entity管理：
	1. 实体Entity-ID和实体Entity的映射关系
	2. 网络层和实体Entity的映射关系
	3. 服务service和实体Entity的关系（举个例子：LOL中温泉是个实体，它提供了一种回血和攻击的服务）
	4. 提供注册Entity的能力
Space管理：
	1. 实体Space-ID和实体Space的映射关系
	2. 提供注册Space的能力
AOI管理：
	AOI管理其实被抽象成十字链表算法，维护着Entity进入、离开、移动时的位置和邻居信息，
	它的存在就是实时记录实体Entity的位置信息和实时计算实体Entity的关注区域
	
规划：
A 所有实体的aoi的由服务器控制.
space中统一管理着entity的进入，离开，和移动
entity有哪些:
1. player
2. 小兵
3. 防御塔
B 技能系统是一套独立的系统，普通攻击也算是一种技能，每种有技能的实体都会绑定一棵AI行为树.
```
