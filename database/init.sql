CREATE TABLE IF NOT EXISTS users (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(80) NOT NULL,
  major VARCHAR(120) NOT NULL,
  credit_score INT NOT NULL DEFAULT 80,
  credit_level VARCHAR(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS skills (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  level_score INT NOT NULL,
  campus VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  portfolio VARCHAR(160) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS needs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id BIGINT NOT NULL,
  title VARCHAR(120) NOT NULL,
  category VARCHAR(40) NOT NULL,
  campus VARCHAR(40) NOT NULL,
  expect_time VARCHAR(80) NOT NULL,
  budget_type VARCHAR(40) NOT NULL,
  description TEXT NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'open'
);

-- 求助候选名单：响应者提交能帮的方式和可上门时段，按提交先后排队
CREATE TABLE IF NOT EXISTS need_responses (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  need_id BIGINT NOT NULL,
  responder VARCHAR(80) NOT NULL,
  help_offer TEXT NOT NULL,
  visit_slots VARCHAR(200) NOT NULL,
  status VARCHAR(20) NOT NULL DEFAULT 'waiting',
  submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS appointments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  pair_name VARCHAR(120) NOT NULL,
  exchange_time VARCHAR(80) NOT NULL,
  place VARCHAR(120) NOT NULL,
  status VARCHAR(40) NOT NULL,
  agenda TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  from_user VARCHAR(80) NOT NULL,
  to_user VARCHAR(80) NOT NULL,
  rating INT NOT NULL,
  content TEXT NOT NULL
);

INSERT INTO users(name, major, credit_score, credit_level) VALUES
('林澈', '新闻传播 2023', 91, '黄金导师'),
('孟野', '音乐表演 2022', 88, '白银协作者'),
('周芮', '统计学 2021', 93, '黄金导师');

INSERT INTO skills(user_id, title, category, level_score, campus, description, portfolio) VALUES
(1, '毕业照人像摄影', '摄影', 92, '东校区', '提供构图、修图和毕业季跟拍，可交换吉他入门课。', '12组校园人像作品'),
(2, '民谣吉他陪练', '乐器', 81, '西校区', '节奏型、弹唱和舞台经验分享，想找人拍宣传照。', '校园音乐节演出视频'),
(3, 'Python 数据分析', '编程', 88, '中心校区', 'pandas、可视化、论文数据清洗辅导。', '3份课程项目证书');

INSERT INTO needs(user_id, title, category, campus, expect_time, budget_type, description, status) VALUES
(2, '找人帮忙拍乐队宣传照', '摄影', '西校区', '本周六上午', '技能交换', '可交换 3 次吉他课，希望会调色和室外构图。', 'open'),
(1, '求教 Python 数据分析', '编程', '中心校区', '周二晚', '小额报酬', '论文问卷数据需要清洗和画图，最好有 pandas 经验。', 'open'),
(1, '想学吉他扫弦入门', '乐器', '东校区', '周三晚', '技能交换', '用摄影课交换吉他基础，希望同校区或线上。', 'open');

INSERT INTO need_responses(need_id, responder, help_offer, visit_slots, status) VALUES
(1, '周芮', '可带全画幅相机和反光板上门，包调色修 9 图。', '周六上午,周六下午', 'interviewing'),
(1, '林澈', '提供毕业照同款人像拍摄，可交换吉他入门课。', '周六上午', 'waiting'),
(1, '许安', '会基础构图和 Lightroom 调色，想积累乐队人像作品。', '周六下午,周日全天', 'waiting'),
(2, '周芮', 'pandas 清洗加可视化一次搞定，可提供论文级代码模板。', '周二晚', 'interviewing'),
(2, '孟野', '会 matplotlib 和问卷统计，可帮忙跑数据。', '周二晚,周三晚', 'waiting'),
(3, '孟野', '民谣扫弦四课时入门，可交换人像摄影课。', '周三晚,周六上午', 'interviewing'),
(3, '周芮', '会基础弹唱，可以一起练习互相纠错。', '周三晚', 'waiting');
